import { afterEach, describe, expect, it, vi } from 'vitest'
import { apiRequest, UnauthenticatedError } from './client'

function problem(status: number, title: string, code?: string) {
  return new Response(JSON.stringify({ title, code }), {
    status,
    headers: { 'Content-Type': 'application/problem+json' },
  })
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('apiRequest', () => {
  it('preserves API status and code from problem details', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(problem(409, 'Already exists', 'HANDLE_ALREADY_EXISTS')))

    await expect(apiRequest('/me/profile', {}, false)).rejects.toMatchObject({
      status: 409,
      code: 'HANDLE_ALREADY_EXISTS',
      message: 'Already exists',
    })
  })

  it('refreshes once and retries the original request', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(problem(401, 'Expired'))
      .mockResolvedValueOnce(new Response('', { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ id: 'user-1', type: 'user' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }))
    vi.stubGlobal('fetch', fetchMock)

    await expect(apiRequest('/auth/me')).resolves.toEqual({ id: 'user-1', type: 'user' })
    expect(fetchMock).toHaveBeenCalledTimes(3)
    expect((fetchMock.mock.calls[1]![0] as Request).url).toBe('https://api.test/auth/users/refresh-tokens')
  })

  it('shares one refresh across concurrent unauthorized requests', async () => {
    let resolveRefresh!: (response: Response) => void
    const fetchMock = vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/users/refresh-tokens')) {
        return new Promise<Response>((resolve) => { resolveRefresh = resolve })
      }

      const callsForUrl = fetchMock.mock.calls.filter(([candidate]) => (
        (candidate as Request).url === request.url
      )).length
      if (callsForUrl === 1) return problem(401, 'Expired')

      return new Response(JSON.stringify({ ok: true }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    })
    vi.stubGlobal('fetch', fetchMock)

    const first = apiRequest<{ ok: boolean }>('/first')
    const second = apiRequest<{ ok: boolean }>('/second')

    await vi.waitFor(() => {
      expect(fetchMock.mock.calls.filter(([request]) => (
        (request as Request).url.endsWith('/auth/users/refresh-tokens')
      ))).toHaveLength(1)
    })
    resolveRefresh(new Response('', { status: 200, headers: { 'Content-Length': '0' } }))

    await expect(Promise.all([first, second])).resolves.toEqual([{ ok: true }, { ok: true }])
    expect(fetchMock.mock.calls.filter(([request]) => (
      (request as Request).url.endsWith('/auth/users/refresh-tokens')
    ))).toHaveLength(1)
  })

  it('treats a missing refresh cookie as an unauthenticated session', async () => {
    vi.stubGlobal('fetch', vi.fn()
      .mockResolvedValueOnce(problem(401, 'Expired'))
      .mockResolvedValueOnce(problem(400, 'Validation error')))

    await expect(apiRequest('/auth/me')).rejects.toBeInstanceOf(UnauthenticatedError)
  })

  it('preserves an unexpected refresh service failure', async () => {
    vi.stubGlobal('fetch', vi.fn()
      .mockResolvedValueOnce(problem(401, 'Expired'))
      .mockResolvedValueOnce(problem(503, 'Identity unavailable')))

    await expect(apiRequest('/auth/me')).rejects.toMatchObject({
      status: 503,
      message: 'Identity unavailable',
    })
  })
})
