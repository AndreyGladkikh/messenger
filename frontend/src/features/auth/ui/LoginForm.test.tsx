import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { screen } from '@testing-library/dom'
import { render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, expect, it, vi } from 'vitest'
import { LoginForm } from './LoginForm'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

it('signs in through POST /auth/users/login and replaces the page with home', async () => {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify({ accessToken: 'a', refreshToken: 'r' }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  }))
  vi.stubGlobal('fetch', fetchMock)
  const router = createMemoryRouter([
    { path: '/auth/login', Component: LoginForm },
    { path: '/', element: <div>Home gate</div> },
  ], { initialEntries: ['/auth/login'] })
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  queryClient.setQueryData(['auth', 'current-user'], { id: 'previous-user', type: 'user' })
  queryClient.setQueryData(['profile', 'previous-user'], { displayName: 'Previous', handle: 'previous' })
  queryClient.setQueryData(['chats'], { peers: [{ id: 'previous-contact' }] })
  render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
  const user = userEvent.setup()

  await user.type(screen.getByLabelText('Login'), 'alice')
  await user.type(screen.getByLabelText('Password'), 'secret')
  await user.click(screen.getByRole('button', { name: 'Sign in' }))

  expect(await screen.findByText('Home gate')).toBeTruthy()
  expect(router.state.location.pathname).toBe('/')
  expect(router.state.historyAction).toBe('REPLACE')
  expect(queryClient.getQueryData(['auth', 'current-user'])).toBeUndefined()
  expect(queryClient.getQueryData(['profile', 'previous-user'])).toBeUndefined()
  expect(queryClient.getQueryData(['chats'])).toBeUndefined()
  const request = fetchMock.mock.calls[0]![0] as Request
  expect(request.method).toBe('POST')
  expect(request.url).toBe('https://api.test/auth/users/login')
})
