import { expect, test, type Page, type Route } from '@playwright/test'

const accountId = 'credentials-form-user'
const queryClientModule = '/src/shared/lib/query-client.ts'
const privateQueryKey = ['credentials-form-test', 'private']
const credentials = { login: ' anya ', password: ' p ' }

const forms = [
  {
    name: 'sign in',
    path: '/auth/login',
    heading: 'Welcome back.',
    submit: 'Sign In',
    pending: 'Signing in…',
    endpoint: '/auth/users/login',
    successStatus: 200,
    autocomplete: 'current-password',
    alternative: 'Sign up',
    alternativePath: '/auth/register',
    hasProfile: true,
    failureStatus: 401,
    failureMessage: 'Incorrect login or password. Please try again.',
    genericStatus: 503,
    genericMessage: 'Could not sign in. Please try again.',
  },
  {
    name: 'sign up',
    path: '/auth/register',
    heading: 'Create your account.',
    submit: 'Sign Up',
    pending: 'Creating account…',
    endpoint: '/auth/users',
    successStatus: 201,
    autocomplete: 'new-password',
    alternative: 'Sign in',
    alternativePath: '/auth/login',
    hasProfile: false,
    failureStatus: 409,
    failureMessage: 'This login is already taken. Try another one.',
    genericStatus: 401,
    genericMessage: 'Could not create your account. Please try again.',
  },
] as const

function problem(route: Route, status: number) {
  return route.fulfill({
    status,
    contentType: 'application/problem+json',
    body: JSON.stringify({ title: 'Request failed', detail: 'Internal service detail' }),
  })
}

async function mockSession(page: Page, hasProfile: boolean, signedIn = false) {
  const session = { signedIn, refreshRequests: 0 }

  await page.route('**/api/auth/me', (route) => session.signedIn
    ? route.fulfill({ json: { id: accountId, type: 'user' } })
    : problem(route, 401))
  await page.route('**/api/me/profile', (route) => hasProfile
    ? route.fulfill({ json: { displayName: 'Anya', handle: 'anya' } })
    : problem(route, 404))
  await page.route('**/api/auth/users/refresh-tokens', (route) => {
    session.refreshRequests += 1
    return problem(route, 401)
  })

  return session
}

async function fillCredentials(page: Page, password = credentials.password) {
  await page.getByLabel('Login', { exact: true }).fill(credentials.login)
  await page.getByLabel('Password', { exact: true }).fill(password)
}

for (const form of forms) {
  test(`${form.name} has required accessible credentials, password visibility and the alternate route`, async ({ page }) => {
    await mockSession(page, form.hasProfile)
    let submissions = 0
    await page.route(`**/api${form.endpoint}`, (route) => {
      submissions += 1
      return problem(route, 503)
    })
    await page.goto(form.path)

    await expect(page.getByRole('heading', { name: form.heading, level: 1 })).toBeVisible()
    const login = page.getByLabel('Login', { exact: true })
    const password = page.getByLabel('Password', { exact: true })
    const submit = page.getByRole('button', { name: form.submit, exact: true })

    await expect(login).toHaveAttribute('autocomplete', 'username')
    if (form.name === 'sign up') {
      await expect(login).toHaveAccessibleDescription('Use this login whenever you sign in.')
    }
    await expect(password).toHaveAttribute('autocomplete', form.autocomplete)
    await submit.click()
    await expect(login).toBeFocused()
    expect(await login.evaluate((input: HTMLInputElement) => input.validity.valueMissing)).toBe(true)
    await login.fill('anya')
    await submit.click()
    await expect(password).toBeFocused()
    expect(await password.evaluate((input: HTMLInputElement) => input.validity.valueMissing)).toBe(true)
    expect(submissions).toBe(0)

    await password.fill('x')
    expect(await password.evaluate((input: HTMLInputElement) => input.validity.valid)).toBe(true)
    await expect(password).toHaveAttribute('type', 'password')
    await page.getByRole('button', { name: 'Show password', exact: true }).click()
    await expect(password).toHaveAttribute('type', 'text')
    await expect(password).toHaveValue('x')
    await page.getByRole('button', { name: 'Hide password', exact: true }).click()
    await expect(password).toHaveAttribute('type', 'password')
    expect(submissions).toBe(0)

    await page.getByRole('link', { name: form.alternative, exact: true }).click()
    await expect(page).toHaveURL(form.alternativePath)
    await expect(page.getByLabel('Password', { exact: true })).toHaveValue('')
  })

  test(`${form.name} sends unchanged credentials and cookies once, clears private cache and replaces the route`, async ({ page, context, baseURL }) => {
    const session = await mockSession(page, form.hasProfile)
    await context.addCookies([{
      name: 'refresh_token',
      value: 'previous-refresh-token',
      url: baseURL!,
      httpOnly: true,
      sameSite: 'Lax',
    }])
    let submissions = 0
    let completeSubmission!: () => void
    const pendingSubmission = new Promise<void>((resolve) => { completeSubmission = resolve })
    await page.route(`**/api${form.endpoint}`, async (route) => {
      submissions += 1
      await pendingSubmission
      session.signedIn = true
      await route.fulfill({
        status: form.successStatus,
        headers: { 'Set-Cookie': 'refresh_token=new-refresh-token; Path=/; HttpOnly; SameSite=Lax' },
      })
    })

    await page.goto(form.path)
    await fillCredentials(page)
    const historyLength = await page.evaluate(() => history.length)
    await page.evaluate(async ({ moduleUrl, queryKey }) => {
      const { default: queryClient } = await import(moduleUrl)
      queryClient.setQueryData(queryKey, { privateMessage: 'Previous account data' })
      queryClient.setQueryData(['profile', 'previous-account'], { displayName: 'Previous account' })
    }, { moduleUrl: queryClientModule, queryKey: privateQueryKey })

    const requestPromise = page.waitForRequest(`**/api${form.endpoint}`)
    await page.getByRole('button', { name: form.submit, exact: true }).click()
    const request = await requestPromise
    expect(request.method()).toBe('POST')
    expect(request.postDataJSON()).toEqual(credentials)
    expect(await request.headerValue('content-type')).toBe('application/json')
    expect(await request.headerValue('cookie')).toContain('refresh_token=previous-refresh-token')
    await expect(page.getByRole('button', { name: form.pending, exact: true })).toBeDisabled()
    await page.keyboard.press('Enter')
    await page.keyboard.press('Enter')
    expect(submissions).toBe(1)
    await expect(page).toHaveURL(form.path)
    completeSubmission()

    await expect(page).toHaveURL(form.hasProfile ? '/' : '/onboarding/profile')
    if (form.hasProfile) await expect(page.getByRole('button', { name: 'Account menu' })).toBeVisible()
    expect(await page.evaluate(() => history.length)).toBe(historyLength)
    expect((await context.cookies()).find((cookie) => cookie.name === 'refresh_token')?.value).toBe('new-refresh-token')
    const cachedData = await page.evaluate(async ({ moduleUrl, queryKey }) => {
      const { default: queryClient } = await import(moduleUrl)
      return {
        privateQuery: queryClient.getQueryData(queryKey),
        previousProfile: queryClient.getQueryData(['profile', 'previous-account']),
        account: queryClient.getQueryData(['auth', 'current-user']),
      }
    }, { moduleUrl: queryClientModule, queryKey: privateQueryKey })
    expect(cachedData).toEqual({
      privateQuery: undefined,
      previousProfile: undefined,
      account: { id: accountId, type: 'user' },
    })
    expect(submissions).toBe(1)
  })

  test(`${form.name} explains credential failures without refresh or retry and accepts a corrected submission`, async ({ page }) => {
    const session = await mockSession(page, form.hasProfile)
    let submissions = 0
    await page.route(`**/api${form.endpoint}`, (route) => {
      submissions += 1
      if (submissions === 1) return problem(route, form.failureStatus)
      session.signedIn = true
      return route.fulfill({ status: form.successStatus })
    })

    await page.goto(form.path)
    await fillCredentials(page)
    const initialRefreshRequests = session.refreshRequests
    await page.getByRole('button', { name: form.submit, exact: true }).click()

    await expect(page.getByRole('alert')).toHaveText(form.failureMessage)
    await expect(page.getByRole('button', { name: form.submit, exact: true })).toBeEnabled()
    await expect(page).toHaveURL(form.path)
    await expect(page.getByLabel('Login', { exact: true })).toHaveValue(credentials.login)
    expect(submissions).toBe(1)
    expect(session.refreshRequests).toBe(initialRefreshRequests)

    await fillCredentials(page, 'corrected-password')
    await page.getByRole('button', { name: form.submit, exact: true }).click()
    await expect(page).toHaveURL(form.hasProfile ? '/' : '/onboarding/profile')
    expect(submissions).toBe(2)
  })

  test(`${form.name} handles other failures without leaking server details or refreshing credentials`, async ({ page }) => {
    const session = await mockSession(page, form.hasProfile)
    let submissions = 0
    await page.route(`**/api${form.endpoint}`, (route) => {
      submissions += 1
      return problem(route, form.genericStatus)
    })

    await page.goto(form.path)
    await fillCredentials(page)
    const initialRefreshRequests = session.refreshRequests
    await page.getByRole('button', { name: form.submit, exact: true }).click()

    await expect(page.getByRole('alert')).toHaveText(form.genericMessage)
    await expect(page.getByRole('button', { name: form.submit, exact: true })).toBeEnabled()
    await expect(page).toHaveURL(form.path)
    expect(submissions).toBe(1)
    expect(session.refreshRequests).toBe(initialRefreshRequests)
  })

  for (const hasProfile of [true, false]) {
    test(`${form.name} redirects an existing session ${hasProfile ? 'home' : 'to profile onboarding'}`, async ({ page }) => {
      await mockSession(page, hasProfile, true)
      await page.goto(form.path)

      await expect(page).toHaveURL(hasProfile ? '/' : '/onboarding/profile')
      await expect(page.getByRole('heading', { name: form.heading, level: 1 })).not.toBeVisible()
    })
  }
}
