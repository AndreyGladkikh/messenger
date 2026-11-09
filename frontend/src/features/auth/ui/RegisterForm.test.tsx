import { screen } from '@testing-library/dom'
import { render } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, expect, it, vi } from 'vitest'
import { RegisterForm } from './RegisterForm'

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

it('registers through POST /auth/users and replaces the page with onboarding', async () => {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify({ accessToken: 'a', refreshToken: 'r' }), {
    status: 201,
    headers: { 'Content-Type': 'application/json' },
  }))
  vi.stubGlobal('fetch', fetchMock)
  const router = createMemoryRouter([
    { path: '/auth/register', Component: RegisterForm },
    { path: '/onboarding/profile', element: <div>Onboarding</div> },
  ], { initialEntries: ['/auth/register'] })
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  queryClient.setQueryData(['chats'], { peers: [{ id: 'previous-contact' }] })
  render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  )
  const user = userEvent.setup()

  await user.type(screen.getByLabelText('Login'), 'alice')
  await user.type(screen.getByLabelText('Password'), 'secret')
  await user.click(screen.getByRole('button', { name: 'Sign up' }))

  expect(await screen.findByText('Onboarding')).toBeTruthy()
  expect(router.state.location.pathname).toBe('/onboarding/profile')
  expect(router.state.historyAction).toBe('REPLACE')
  expect(queryClient.getQueryData(['chats'])).toBeUndefined()
  const request = fetchMock.mock.calls[0]![0] as Request
  expect(request.method).toBe('POST')
  expect(request.url).toBe('https://api.test/auth/users')
})
