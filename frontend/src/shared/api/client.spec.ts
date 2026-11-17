import { expect, test, type Request } from '@playwright/test'

const clientModule = '/src/shared/api/client.ts'
const avatarBytes = [137, 80, 78, 71, 13, 10, 26, 10, 0, 255, 128]

async function readMultipart(request: Request) {
  const contentType = await request.headerValue('content-type')
  expect(contentType).toMatch(/^multipart\/form-data; boundary=.+/)
  const body = request.postDataBuffer()
  expect(body).not.toBeNull()
  const form = await new Response(new Uint8Array(body!), {
    headers: { 'Content-Type': contentType! },
  }).formData()
  const avatar = form.get('avatar')
  if (!avatar || typeof avatar === 'string') throw new Error('Expected an avatar file')

  return {
    displayName: form.get('displayName'),
    handle: form.get('handle'),
    avatar: {
      name: avatar.name,
      type: avatar.type,
      bytes: [...new Uint8Array(await avatar.arrayBuffer())],
    },
  }
}

test.beforeEach(async ({ page }) => {
  // Load the transport in a browser without starting the application's session queries.
  await page.route('**/transport-test', (route) => route.fulfill({
    contentType: 'text/html',
    body: '<!doctype html><html><head><title>API transport test</title></head><body></body></html>',
  }))
  await page.goto('/transport-test')
})

test('multipart requests preserve fields and file bytes with a browser-generated boundary and session cookies', async ({ page, context, baseURL }) => {
  await context.addCookies([
    { name: 'access_token', value: 'access-token', url: baseURL!, sameSite: 'Lax' },
    { name: 'refresh_token', value: 'refresh-token', url: baseURL!, httpOnly: true, sameSite: 'Lax' },
  ])
  await page.route('**/api/transport/upload', (route) => route.fulfill({ status: 201, json: { created: true } }))
  const requestPromise = page.waitForRequest('**/api/transport/upload')

  const result = await page.evaluate(async ({ moduleUrl, bytes }) => {
    const { apiRequest } = await import(moduleUrl)
    const body = new FormData()
    body.append('displayName', 'Аня Sokolova')
    body.append('handle', 'anya')
    body.append('avatar', new File([new Uint8Array(bytes)], 'avatar.png', { type: 'image/png' }))
    return apiRequest('/transport/upload', {
      method: 'POST',
      // A caller's stale JSON header must not override the browser's multipart boundary.
      headers: new Headers({ 'Content-Type': 'application/json', 'X-Request-ID': 'upload-id' }),
      body,
    })
  }, { moduleUrl: clientModule, bytes: avatarBytes })

  const request = await requestPromise
  expect(result).toEqual({ created: true })
  expect(request.method()).toBe('POST')
  expect(await readMultipart(request)).toEqual({
    displayName: 'Аня Sokolova',
    handle: 'anya',
    avatar: { name: 'avatar.png', type: 'image/png', bytes: avatarBytes },
  })
  expect(await request.headerValue('x-request-id')).toBe('upload-id')
  expect(await request.headerValue('authorization')).toBe('Bearer access-token')
  expect(await request.headerValue('cookie')).toContain('refresh_token=refresh-token')
})

test('JSON defaults and explicit Headers overrides preserve the requested payload and headers', async ({ page }) => {
  await page.route('**/api/transport/**', (route) => route.fulfill({ status: 204 }))
  const defaultRequestPromise = page.waitForRequest('**/api/transport/json')
  const customRequestPromise = page.waitForRequest('**/api/transport/text')

  await page.evaluate(async (moduleUrl) => {
    const { apiRequest } = await import(moduleUrl)
    await apiRequest('/transport/json', {
      method: 'POST',
      headers: new Headers({ 'X-Request-ID': 'json-id' }),
      body: { displayName: 'Anya' },
    })
    await apiRequest('/transport/text', {
      method: 'POST',
      headers: new Headers({ 'Content-Type': 'text/plain', 'X-Request-ID': 'text-id' }),
      body: 'unchanged text payload',
    })
  }, clientModule)

  const jsonRequest = await defaultRequestPromise
  expect(await jsonRequest.headerValue('content-type')).toBe('application/json')
  expect(await jsonRequest.headerValue('x-request-id')).toBe('json-id')
  expect(jsonRequest.postDataJSON()).toEqual({ displayName: 'Anya' })
  const textRequest = await customRequestPromise
  expect(await textRequest.headerValue('content-type')).toBe('text/plain')
  expect(await textRequest.headerValue('x-request-id')).toBe('text-id')
  expect(textRequest.postData()).toBe('unchanged text payload')
})

test('an unauthorized multipart request refreshes the token and resends the same fields and file', async ({ page, context, baseURL }) => {
  await context.addCookies([
    { name: 'access_token', value: 'expired-token', url: baseURL!, sameSite: 'Lax' },
    { name: 'refresh_token', value: 'refresh-token', url: baseURL!, httpOnly: true, sameSite: 'Lax' },
  ])
  const uploads: Request[] = []
  let refreshRequests = 0
  await page.route('**/api/transport/upload', (route) => {
    uploads.push(route.request())
    return uploads.length === 1
      ? route.fulfill({ status: 401, contentType: 'application/problem+json', json: { title: 'Token expired' } })
      : route.fulfill({ status: 201, json: { created: true } })
  })
  await page.route('**/api/auth/users/refresh-tokens', (route) => {
    refreshRequests += 1
    return route.fulfill({
      status: 204,
      headers: { 'Set-Cookie': 'access_token=fresh-token; Path=/; SameSite=Lax' },
    })
  })

  const result = await page.evaluate(async ({ moduleUrl, bytes }) => {
    const { apiRequest } = await import(moduleUrl)
    const body = new FormData()
    body.append('displayName', 'Anya')
    body.append('handle', 'anya')
    body.append('avatar', new File([new Uint8Array(bytes)], 'avatar.png', { type: 'image/png' }))
    return apiRequest('/transport/upload', { method: 'POST', body })
  }, { moduleUrl: clientModule, bytes: avatarBytes })

  expect(result).toEqual({ created: true })
  expect(refreshRequests).toBe(1)
  expect(uploads).toHaveLength(2)
  expect(await readMultipart(uploads[0])).toEqual({
    displayName: 'Anya',
    handle: 'anya',
    avatar: { name: 'avatar.png', type: 'image/png', bytes: avatarBytes },
  })
  expect(await readMultipart(uploads[1])).toEqual(await readMultipart(uploads[0]))
  expect(await uploads[0].headerValue('authorization')).toBe('Bearer expired-token')
  expect(await uploads[1].headerValue('authorization')).toBe('Bearer fresh-token')
  expect(await uploads[1].headerValue('cookie')).toContain('refresh_token=refresh-token')
})
