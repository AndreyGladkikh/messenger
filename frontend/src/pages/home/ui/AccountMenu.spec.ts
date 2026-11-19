import { expect, test, type Page, type Route } from '@playwright/test'
import type { Profile } from '@/entities/profile'

const accountId = 'account-menu-user'
const queryClientModule = '/src/shared/lib/query-client.ts'
const privateQueryKey = ['account-menu-test', 'private']

function problem(route: Route, status: number) {
  return route.fulfill({
    status,
    contentType: 'application/problem+json',
    body: JSON.stringify({ title: status === 401 ? 'Unauthorized' : 'Service unavailable' }),
  })
}

async function mockSession(page: Page, profile: Profile = { displayName: 'Аня Соколова', handle: 'anya' }) {
  const session = { signedIn: true, refreshRequests: 0 }

  await page.route('**/api/auth/me', (route) => session.signedIn
    ? route.fulfill({ json: { id: accountId, type: 'user' } })
    : problem(route, 401))
  await page.route('**/api/me/profile', (route) => route.fulfill({
    json: profile,
  }))
  await page.route('**/api/auth/users/refresh-tokens', (route) => {
    session.refreshRequests += 1
    return problem(route, 401)
  })

  return session
}

test('avatar menu shows the account and an icon, supports keyboard use and dismissal', async ({ page }) => {
  await mockSession(page)
  await page.goto('/')

  const trigger = page.getByRole('button', { name: 'Account menu' })
  const menu = page.getByRole('menu')
  const signOut = page.getByRole('menuitem', { name: 'Sign out', exact: true })

  await expect(trigger.getByText('АС', { exact: true })).toBeVisible()
  await expect(menu).not.toBeVisible()
  await trigger.click()
  await expect(trigger).toHaveAttribute('aria-expanded', 'true')
  await expect(menu.getByText('Аня Соколова', { exact: true })).toBeVisible()
  await expect(menu.getByText('@anya', { exact: true })).toBeVisible()
  await expect(signOut.locator('svg[aria-hidden="true"]')).toBeVisible()

  await page.keyboard.press('Escape')
  await expect(menu).not.toBeVisible()
  await expect(trigger).toBeFocused()

  await trigger.press('ArrowDown')
  await expect(menu).toBeVisible()
  await expect(signOut).toBeFocused()
  await page.keyboard.press('Escape')
  await expect(trigger).toBeFocused()

  await trigger.press('Enter')
  await expect(menu).toBeVisible()
  await page.mouse.click(1000, 100)
  await expect(menu).not.toBeVisible()
  await expect(trigger).toHaveAttribute('aria-expanded', 'false')
})

test('the account menu renders the profile avatar URL directly and refreshes it after reload', async ({ page, baseURL }) => {
  let profileReads = 0
  let mediaUrlReads = 0
  await page.route('**/api/auth/me', (route) => route.fulfill({ json: { id: accountId, type: 'user' } }))
  await page.route('**/api/me/profile', (route) => {
    profileReads += 1
    return route.fulfill({ json: {
      displayName: 'Mira Chen',
      handle: 'mira_chen',
      avatarUrl: `${baseURL}/saved-avatar.png?version=${profileReads}`,
    } })
  })
  await page.route('**/api/media/*/url', (route) => {
    mediaUrlReads += 1
    return route.fulfill({ status: 500 })
  })
  await page.route('**/saved-avatar.png?version=*', (route) => route.fulfill({
    contentType: 'image/png',
    body: Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGP4z8DwHwAFAAH/iZk9HQAAAABJRU5ErkJggg==', 'base64'),
  }))
  await page.goto('/')

  const trigger = page.getByRole('button', { name: 'Account menu' })
  const avatar = trigger.getByRole('img', { name: 'Mira Chen profile photo' })
  await expect(avatar).toHaveAttribute('src', `${baseURL}/saved-avatar.png?version=1`)
  await expect(avatar).toHaveJSProperty('naturalWidth', 1)
  await trigger.click()
  const menu = page.getByRole('menu')
  await expect(menu.getByText('Mira Chen', { exact: true })).toBeVisible()
  await expect(menu.getByText('@mira_chen', { exact: true })).toBeVisible()

  await page.reload()
  await expect(avatar).toHaveAttribute('src', `${baseURL}/saved-avatar.png?version=2`)
  await expect(avatar).toHaveJSProperty('naturalWidth', 1)
  expect(profileReads).toBe(2)
  expect(mediaUrlReads).toBe(0)
})

test('an unavailable profile image falls back to the saved profile initials', async ({ page }) => {
  const avatarUrl = 'https://cdn.example.test/avatar.png'
  await mockSession(page, { displayName: 'Mira Chen', handle: 'mira_chen', avatarUrl })
  await page.route(avatarUrl, (route) => route.fulfill({ status: 404 }))
  const avatarRequest = page.waitForRequest(avatarUrl)
  await page.goto('/')
  expect((await avatarRequest).url()).toBe(avatarUrl)

  const trigger = page.getByRole('button', { name: 'Account menu' })
  await expect(trigger.getByRole('img')).toHaveCount(0)
  await expect(trigger.getByText('MC', { exact: true })).toBeVisible()
  await trigger.click()
  await expect(page.getByRole('menu').getByText('Mira Chen', { exact: true })).toBeVisible()
})

test('a profile without an avatar or username shows initials and omits the username line', async ({ page }) => {
  await mockSession(page, { displayName: 'Mira Chen' })

  await page.goto('/')

  const trigger = page.getByRole('button', { name: 'Account menu' })
  await expect(trigger.getByRole('img')).toHaveCount(0)
  await expect(trigger.getByText('MC', { exact: true })).toBeVisible()
  await trigger.click()
  await expect(page.getByRole('menu').getByText('Mira Chen', { exact: true })).toBeVisible()
  await expect(page.getByRole('menu').getByText(/^@/)).toHaveCount(0)
})

test('sign out sends the session cookie once, waits, clears private cache and replaces the page', async ({ page, context, baseURL }) => {
  const session = await mockSession(page)
  await context.addCookies([{
    name: 'refresh_token',
    value: 'test-refresh-token',
    url: baseURL!,
    httpOnly: true,
    sameSite: 'Lax',
  }])

  let logoutRequests = 0
  let completeLogout!: () => void
  const pendingLogout = new Promise<void>((resolve) => { completeLogout = resolve })
  await page.route('**/api/auth/users/logout', async (route) => {
    logoutRequests += 1
    await pendingLogout
    session.signedIn = false
    await route.fulfill({
      status: 204,
      headers: { 'Set-Cookie': 'refresh_token=; Max-Age=0; Path=/; HttpOnly; SameSite=Lax' },
    })
  })

  await page.goto('/')
  await expect(page.getByRole('button', { name: 'Account menu' })).toBeVisible()
  const historyLength = await page.evaluate(() => history.length)
  await page.evaluate(async ({ moduleUrl, queryKey }) => {
    const { default: queryClient } = await import(moduleUrl)
    queryClient.setQueryData(queryKey, { privateMessage: 'Only for the current account' })
  }, { moduleUrl: queryClientModule, queryKey: privateQueryKey })

  await page.getByRole('button', { name: 'Account menu' }).click()
  const logoutRequestPromise = page.waitForRequest('**/api/auth/users/logout')
  await page.getByRole('menuitem', { name: 'Sign out', exact: true }).click()
  const logoutRequest = await logoutRequestPromise

  expect(logoutRequest.method()).toBe('POST')
  expect(logoutRequest.postData()).toBeNull()
  expect(await logoutRequest.headerValue('cookie')).toContain('refresh_token=test-refresh-token')
  const pendingItem = page.getByRole('menuitem', { name: 'Signing out…' })
  await expect(pendingItem).toBeVisible()
  await expect(pendingItem).toBeDisabled()
  await page.keyboard.press('Enter')
  await page.keyboard.press('Enter')
  expect(logoutRequests).toBe(1)
  await expect(page).toHaveURL('/')

  completeLogout()

  await expect(page).toHaveURL('/auth/login', { timeout: 10_000 })
  await expect(page.getByRole('button', { name: 'Sign In', exact: true })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Account menu' })).not.toBeVisible()
  expect(await page.evaluate(() => history.length)).toBe(historyLength)
  expect((await context.cookies()).some((cookie) => cookie.name === 'refresh_token')).toBe(false)
  const privateData = await page.evaluate(async ({ moduleUrl, queryKey, id }) => {
    const { default: queryClient } = await import(moduleUrl)
    return {
      privateQuery: queryClient.getQueryData(queryKey),
      profile: queryClient.getQueryData(['profile', id]),
    }
  }, { moduleUrl: queryClientModule, queryKey: privateQueryKey, id: accountId })
  expect(privateData).toEqual({ privateQuery: undefined, profile: undefined })
  expect(logoutRequests).toBe(1)
})

test('a failed sign out keeps the session and allows a successful retry', async ({ page }) => {
  const session = await mockSession(page)
  let logoutRequests = 0
  await page.route('**/api/auth/users/logout', (route) => {
    logoutRequests += 1
    if (logoutRequests === 1) return problem(route, 503)
    session.signedIn = false
    return route.fulfill({ status: 204 })
  })

  await page.goto('/')
  await page.getByRole('button', { name: 'Account menu' }).click()
  await page.getByRole('menuitem', { name: 'Sign out', exact: true }).click()

  await expect(page.getByRole('alert')).toHaveText('Could not sign out. Please try again.')
  await expect(page).toHaveURL('/')
  await expect(page.getByRole('menuitem', { name: 'Sign out', exact: true })).toBeEnabled()
  expect(logoutRequests).toBe(1)
  const profile = await page.evaluate(async ({ moduleUrl, id }) => {
    const { default: queryClient } = await import(moduleUrl)
    return queryClient.getQueryData(['profile', id])
  }, { moduleUrl: queryClientModule, id: accountId })
  expect(profile).toEqual({ displayName: 'Аня Соколова', handle: 'anya' })

  await page.getByRole('menuitem', { name: 'Sign out', exact: true }).click()
  await expect(page).toHaveURL('/auth/login', { timeout: 10_000 })
  await expect(page.getByRole('button', { name: 'Sign In', exact: true })).toBeVisible()
  expect(logoutRequests).toBe(2)
})

test('an unauthorized sign out reports failure without refreshing or retrying the request', async ({ page }) => {
  const session = await mockSession(page)
  let logoutRequests = 0
  await page.route('**/api/auth/users/logout', (route) => {
    logoutRequests += 1
    return problem(route, 401)
  })

  await page.goto('/')
  await page.getByRole('button', { name: 'Account menu' }).click()
  await page.getByRole('menuitem', { name: 'Sign out', exact: true }).click()

  await expect(page.getByRole('alert')).toHaveText('Could not sign out. Please try again.')
  await expect(page).toHaveURL('/')
  await expect(page.getByRole('menuitem', { name: 'Sign out', exact: true })).toBeEnabled()
  expect(logoutRequests).toBe(1)
  expect(session.refreshRequests).toBe(0)
})
