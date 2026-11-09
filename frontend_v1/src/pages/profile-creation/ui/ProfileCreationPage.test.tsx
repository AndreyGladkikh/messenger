import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { screen } from '@testing-library/dom'
import { render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ProfileCreationPage } from './ProfileCreationPage'

const account = { id: 'user-1', type: 'user' }

function problem(status: number, title: string, code?: string) {
  return new Response(JSON.stringify({ title, code }), {
    status,
    headers: { 'Content-Type': 'application/problem+json' },
  })
}

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  const router = createMemoryRouter([
    { path: '/', element: <div>Home</div> },
    {
      path: '/onboarding/profile',
      loader: () => account,
      HydrateFallback: () => <div>Loading…</div>,
      Component: ProfileCreationPage,
    },
  ], { initialEntries: ['/onboarding/profile'] })

  render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )

  return { queryClient, router }
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('ProfileCreationPage', () => {
  it('requires nonblank display name and username without sending a request', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    renderPage()

    await userEvent.click(await screen.findByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Display name is required.')).toBeTruthy()
    expect(screen.getByText('Username is required.')).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('validates an optional avatar as an absolute HTTP URL', async () => {
    vi.stubGlobal('fetch', vi.fn())
    renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.type(screen.getByLabelText('Avatar URL (optional)'), 'avatar.png')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Enter an absolute http:// or https:// URL.')).toBeTruthy()
  })

  it('submits a valid avatar URL', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('', {
      status: 201,
      headers: { 'Content-Length': '0' },
    }))
    vi.stubGlobal('fetch', fetchMock)
    renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.type(screen.getByLabelText('Avatar URL (optional)'), 'https://example.com/alice.png')

    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Home')).toBeTruthy()
    const request = fetchMock.mock.calls[0]![0] as Request
    expect(await request.json()).toEqual({
      displayName: 'Alice',
      handle: 'alice',
      avatarUrl: 'https://example.com/alice.png',
    })
  })

  it('blocks duplicate submissions while profile creation is pending', async () => {
    let resolveRequest!: (response: Response) => void
    const fetchMock = vi.fn().mockImplementation(() => new Promise<Response>((resolve) => {
      resolveRequest = resolve
    }))
    vi.stubGlobal('fetch', fetchMock)
    renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    await vi.waitFor(() => {
      const button = screen.getByRole('button', { name: 'Create profile' }) as HTMLButtonElement
      expect(button.disabled).toBe(true)
      expect(button.getAttribute('aria-busy')).toBe('true')
    })
    await user.click(screen.getByRole('button', { name: 'Create profile' }))
    expect(fetchMock).toHaveBeenCalledTimes(1)

    resolveRequest(new Response('', { status: 201, headers: { 'Content-Length': '0' } }))
    expect(await screen.findByText('Home')).toBeTruthy()
  })

  it('blocks duplicate submissions while a conflict is being rechecked', async () => {
    let postCount = 0
    let recheckStarted = false
    let resolveRecheck!: (response: Response) => void
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.method === 'POST') {
        postCount += 1
        return problem(409, 'Already exists', 'HANDLE_ALREADY_EXISTS')
      }
      recheckStarted = true
      return new Promise<Response>((resolve) => { resolveRecheck = resolve })
    }))
    renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    await vi.waitFor(() => expect(recheckStarted).toBe(true))
    const button = screen.getByRole('button', { name: 'Create profile' }) as HTMLButtonElement
    expect(button.disabled).toBe(true)
    expect(button.getAttribute('aria-busy')).toBe('true')
    await user.click(button)
    expect(postCount).toBe(1)

    resolveRecheck(problem(404, 'Not found'))
    expect(await screen.findByText('This username is already taken.')).toBeTruthy()
  })

  it('preserves values and shows a field error when the handle is taken', async () => {
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.method === 'POST') return problem(409, 'Already exists', 'HANDLE_ALREADY_EXISTS')
      return problem(404, 'Not found')
    }))
    renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'taken')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('This username is already taken.')).toBeTruthy()
    expect((screen.getByLabelText('Display Name') as HTMLInputElement).value).toBe('Alice')
    expect((screen.getByLabelText('Username / Handle') as HTMLInputElement).value).toBe('taken')
  })

  it('creates a profile without an avatar and replaces onboarding with home', async () => {
    const fetchMock = vi.fn().mockResolvedValue(new Response('', {
      status: 201,
      headers: { 'Content-Length': '0' },
    }))
    vi.stubGlobal('fetch', fetchMock)
    const { router } = renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), ' Alice ')
    await user.type(screen.getByLabelText('Username / Handle'), ' alice ')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Home')).toBeTruthy()
    expect(router.state.location.pathname).toBe('/')
    expect(router.state.historyAction).toBe('REPLACE')
    const request = fetchMock.mock.calls[0]![0] as Request
    expect(await request.json()).toEqual({ displayName: 'Alice', handle: 'alice' })
  })

  it('recovers an ambiguous conflict when the profile was already created', async () => {
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.method === 'POST') return problem(409, 'Already exists', 'DIRECTORY_ENTRY_ALREADY_EXISTS')
      return new Response(JSON.stringify({ displayName: 'Alice', handle: 'alice' }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      })
    }))
    const { router } = renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Home')).toBeTruthy()
    expect(router.state.location.pathname).toBe('/')
  })

  it('stays on onboarding when a conflict recheck fails', async () => {
    vi.stubGlobal('fetch', vi.fn(async (request: Request) => {
      if (request.method === 'POST') return problem(409, 'Already exists', 'DIRECTORY_ENTRY_ALREADY_EXISTS')
      return problem(503, 'Directory unavailable')
    }))
    const { router } = renderPage()
    const user = userEvent.setup()

    await user.type(await screen.findByLabelText('Display Name'), 'Alice')
    await user.type(screen.getByLabelText('Username / Handle'), 'alice')
    await user.click(screen.getByRole('button', { name: 'Create profile' }))

    expect(await screen.findByText('Directory unavailable')).toBeTruthy()
    expect(router.state.location.pathname).toBe('/onboarding/profile')
    expect((screen.getByLabelText('Username / Handle') as HTMLInputElement).value).toBe('alice')
  })
})
