import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { screen } from '@testing-library/dom'
import { act, render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterContextProvider, RouterProvider } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { queryClientContext } from './context'
import { APP_ROUTES } from './routes'

function json(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

function problem(status: number, title: string) {
  return new Response(JSON.stringify({ title }), {
    status,
    headers: { 'Content-Type': 'application/problem+json' },
  })
}

function renderRoutes(initialEntry: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter(APP_ROUTES, {
    initialEntries: [initialEntry],
    getContext() {
      const context = new RouterContextProvider()
      context.set(queryClientContext, queryClient)
      return context
    },
  })

  render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
  return router
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('profile route gate', () => {
  it.each([
    ['/auth/login', 'Welcome back'],
    ['/auth/register', 'Create your account'],
  ])('allows a fresh visitor to open %s when refresh has no cookie', async (path, heading) => {
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return problem(401, 'Sign in required')
      if (request.url.endsWith('/auth/users/refresh-tokens')) return problem(400, 'Validation error')
      throw new Error(`Unexpected request: ${request.url}`)
    }))
    const router = renderRoutes(path)

    expect(await screen.findByRole('heading', { name: heading })).toBeTruthy()
    expect(router.state.location.pathname).toBe(path)
  })

  it('redirects an unauthenticated home visit to login', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => problem(401, 'Sign in required')))
    const router = renderRoutes('/')

    expect(await screen.findByRole('heading', { name: 'Welcome back' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/auth/login')
  })

  it('redirects an incomplete account before chat or WebSocket can mount', async () => {
    let websocketCount = 0
    class WebSocketStub {
      onopen: (() => void) | null = null
      onmessage: ((event: MessageEvent) => void) | null = null
      constructor() { websocketCount += 1 }
      close() {}
    }
    vi.stubGlobal('WebSocket', WebSocketStub)
    const fetchMock = vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return json({ id: 'user-1', type: 'user' })
      if (request.url.endsWith('/me/profile')) return problem(404, 'Not found')
      throw new Error(`Unexpected messenger request: ${request.url}`)
    })
    vi.stubGlobal('fetch', fetchMock)
    const router = renderRoutes('/')

    expect(await screen.findByRole('heading', { name: 'Create your profile' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/onboarding/profile')
    expect(websocketCount).toBe(0)
    expect(fetchMock.mock.calls.some(([request]) => (request as Request).url.endsWith('/bff/web/home-page'))).toBe(false)
  })

  it.each(['/auth/login', '/auth/register'])('redirects an incomplete account from %s to onboarding', async (path) => {
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return json({ id: 'user-1', type: 'user' })
      if (request.url.endsWith('/me/profile')) return problem(404, 'Not found')
      throw new Error(`Unexpected request: ${request.url}`)
    }))
    const router = renderRoutes(path)

    expect(await screen.findByRole('heading', { name: 'Create your profile' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/onboarding/profile')
  })

  it('keeps messenger locked on a profile lookup failure and retries the gate', async () => {
    let profileAttempts = 0
    let websocketCount = 0
    class WebSocketStub {
      onopen: (() => void) | null = null
      onmessage: ((event: MessageEvent) => void) | null = null
      constructor() { websocketCount += 1 }
      close() {}
    }
    vi.stubGlobal('WebSocket', WebSocketStub)
    const fetchMock = vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return json({ id: 'user-1', type: 'user' })
      if (request.url.endsWith('/me/profile')) {
        profileAttempts += 1
        return profileAttempts === 1
          ? problem(503, 'Directory unavailable')
          : json({ displayName: 'Alice', handle: 'alice' })
      }
      if (request.url.endsWith('/bff/web/home-page')) return json({ peers: [] })
      throw new Error(`Unexpected request: ${request.url}`)
    })
    vi.stubGlobal('fetch', fetchMock)
    renderRoutes('/')

    expect(await screen.findByRole('heading', { name: 'We couldn’t load this page' })).toBeTruthy()
    expect(screen.getByText('Directory unavailable')).toBeTruthy()
    expect(websocketCount).toBe(0)
    expect(fetchMock.mock.calls.some(([request]) => (request as Request).url.endsWith('/bff/web/home-page'))).toBe(false)

    await userEvent.click(screen.getByRole('button', { name: 'Retry' }))

    expect(await screen.findByRole('heading', { name: 'Messages' })).toBeTruthy()
    expect(websocketCount).toBe(1)
  })

  it('allows an account with a profile to enter home', async () => {
    let websocketCount = 0
    class WebSocketStub {
      onopen: (() => void) | null = null
      onmessage: ((event: MessageEvent) => void) | null = null
      constructor() { websocketCount += 1 }
      close() {}
    }
    vi.stubGlobal('WebSocket', WebSocketStub)
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return json({ id: 'user-1', type: 'user' })
      if (request.url.endsWith('/me/profile')) return json({ displayName: 'Alice', handle: 'alice' })
      if (request.url.endsWith('/bff/web/home-page')) return json({ peers: [] })
      throw new Error(`Unexpected request: ${request.url}`)
    }))
    const router = renderRoutes('/')

    expect(await screen.findByRole('heading', { name: 'Messages' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/')
    expect(websocketCount).toBe(1)
  })

  it('unmounts messenger while a new account profile is being verified', async () => {
    let accountId = 'user-1'
    let profileCheckStarted = false
    let resolveProfile!: (response: Response) => void
    let websocketCount = 0
    class WebSocketStub {
      onopen: (() => void) | null = null
      onmessage: ((event: MessageEvent) => void) | null = null
      constructor() { websocketCount += 1 }
      close() { websocketCount -= 1 }
    }
    vi.stubGlobal('WebSocket', WebSocketStub)
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.url.endsWith('/auth/me')) return json({ id: accountId, type: 'user' })
      if (request.url.endsWith('/me/profile')) {
        if (accountId === 'user-1') return json({ displayName: 'Alice', handle: 'alice' })
        if (profileCheckStarted) return problem(404, 'Not found')
        profileCheckStarted = true
        return new Promise<Response>((resolve) => { resolveProfile = resolve })
      }
      if (request.url.endsWith('/bff/web/home-page')) return json({ peers: [] })
      throw new Error(`Unexpected request: ${request.url}`)
    }))
    const router = renderRoutes('/')

    expect(await screen.findByRole('heading', { name: 'Messages' })).toBeTruthy()
    expect(websocketCount).toBe(1)
    accountId = 'user-2'

    let navigation!: Promise<void>
    await act(async () => {
      navigation = router.navigate('/?account=user-2')
      await Promise.resolve()
    })
    await vi.waitFor(() => expect(profileCheckStarted).toBe(true))

    expect(screen.queryByRole('heading', { name: 'Messages' })).toBeNull()
    expect(screen.getByText('Loading…')).toBeTruthy()
    expect(websocketCount).toBe(0)

    await act(async () => {
      resolveProfile(problem(404, 'Not found'))
      await navigation
    })
    expect(await screen.findByRole('heading', { name: 'Create your profile' })).toBeTruthy()
  })

  it.each(['/auth/login', '/auth/register', '/onboarding/profile'])(
    'redirects a completed account from %s to home',
    async (path) => {
      class WebSocketStub {
        onopen: (() => void) | null = null
        onmessage: ((event: MessageEvent) => void) | null = null
        close() {}
      }
      vi.stubGlobal('WebSocket', WebSocketStub)
      vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
        if (request.url.endsWith('/auth/me')) return json({ id: 'user-1', type: 'user' })
        if (request.url.endsWith('/me/profile')) return json({ displayName: 'Alice', handle: 'alice' })
        if (request.url.endsWith('/bff/web/home-page')) return json({ peers: [] })
        throw new Error(`Unexpected request: ${request.url}`)
      }))
      const router = renderRoutes(path)

      expect(await screen.findByRole('heading', { name: 'Messages' })).toBeTruthy()
      expect(router.state.location.pathname).toBe('/')
    },
  )
})
