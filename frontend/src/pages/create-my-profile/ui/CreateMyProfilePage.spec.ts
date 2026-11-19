/// <reference types="node" />
import { Buffer } from 'node:buffer'
import { expect, test, type Page, type Route } from '@playwright/test'

const onboardingPath = '/onboarding/profile'
const accountId = 'profile-onboarding-user'
const profile = { displayName: 'Anya Sokolova', handle: 'anya' }
const avatarImage = {
  name: 'avatar.png',
  mimeType: 'image/png',
  buffer: Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aD1sAAAAASUVORK5CYII=', 'base64'),
}
const storageOrigin = 'http://127.0.0.1:4173'

interface SubmittedProfile {
  displayName: string
  handle?: string
  bio?: string
  birthday?: string
  avatarMediaId?: string
}

interface SavedProfile extends Omit<SubmittedProfile, 'avatarMediaId'> {
  avatarUrl?: string
}

interface MockSession {
  signedIn: boolean
  hasProfile: boolean
  profileStatus?: number
  savedProfile: SavedProfile
  profileReads: number
  submissions: number
  uploadInitializations: number
  storageUploads: number
  uploadCompletions: number
  avatarUrlReads: number
  requestOrder: string[]
  create: (route: Route) => Promise<void>
  initUpload: (route: Route) => Promise<void>
  putUpload: (route: Route) => Promise<void>
  completeUpload: (route: Route) => Promise<void>
}

function savedProfileFrom(route: Route): SavedProfile {
  const { avatarMediaId, ...submitted } = route.request().postDataJSON() as SubmittedProfile
  return {
    ...submitted,
    ...(avatarMediaId ? { avatarUrl: `${storageOrigin}/object-storage/saved-avatar` } : {}),
  }
}

function problem(route: Route, status: number, code?: string) {
  return route.fulfill({
    status,
    contentType: 'application/problem+json',
    body: JSON.stringify({ title: 'Request failed', detail: 'Internal service detail', code }),
  })
}

async function mockSession(page: Page, options: { signedIn?: boolean, hasProfile?: boolean, profileStatus?: number } = {}) {
  const session: MockSession = {
    signedIn: options.signedIn ?? true,
    hasProfile: options.hasProfile ?? false,
    profileStatus: options.profileStatus,
    savedProfile: { ...profile },
    profileReads: 0,
    submissions: 0,
    uploadInitializations: 0,
    storageUploads: 0,
    uploadCompletions: 0,
    avatarUrlReads: 0,
    requestOrder: [],
    create: (route: Route): Promise<void> => {
      session.savedProfile = savedProfileFrom(route)
      session.hasProfile = true
      return route.fulfill({ status: 201 })
    },
    initUpload: (route: Route): Promise<void> => route.fulfill({
      status: 201,
      json: {
        id: `upload-${session.uploadInitializations}`,
        presignedUrl: `${storageOrigin}/object-storage/upload-${session.uploadInitializations}`,
      },
    }),
    putUpload: (route: Route): Promise<void> => route.fulfill({ status: 200 }),
    completeUpload: (route: Route): Promise<void> => route.fulfill({
      json: { mediaId: `media-${session.uploadCompletions}` },
    }),
  }

  await page.route('**/api/auth/me', (route) => session.signedIn
    ? route.fulfill({ json: { id: accountId, type: 'user' } })
    : problem(route, 401))
  await page.route('**/api/auth/users/refresh-tokens', (route) => problem(route, 401))
  await page.route('**/api/me/profile', (route) => {
    if (route.request().method() === 'POST') {
      session.submissions += 1
      return session.create(route)
    }
    session.profileReads += 1
    if (session.profileStatus) return problem(route, session.profileStatus)
    return session.hasProfile ? route.fulfill({ json: session.savedProfile }) : problem(route, 404)
  })
  await page.route('**/api/media/uploads/*/complete', (route) => {
    session.uploadCompletions += 1
    session.requestOrder.push('complete')
    return session.completeUpload(route)
  })
  await page.route('**/api/media/uploads', (route) => {
    session.uploadInitializations += 1
    session.requestOrder.push('init')
    return session.initUpload(route)
  })
  await page.route('**/api/media/*/url', (route) => {
    session.avatarUrlReads += 1
    return route.fulfill({ json: { url: `${storageOrigin}/object-storage/saved-avatar` } })
  })
  await page.route('**/object-storage/*', (route) => {
    if (route.request().method() === 'GET') {
      return route.fulfill({ contentType: avatarImage.mimeType, body: avatarImage.buffer })
    }
    session.storageUploads += 1
    session.requestOrder.push('put')
    return session.putUpload(route)
  })

  return session
}

async function fillProfile(page: Page, displayName = profile.displayName, username = profile.handle) {
  await page.getByLabel('Display Name', { exact: true }).fill(displayName)
  await page.getByLabel('Username (optional)', { exact: true }).fill(username)
}

async function expectMessengerLocked(page: Page) {
  await expect(page.getByRole('button', { name: 'Account menu' })).not.toBeVisible()
  await expect(page.getByRole('textbox', { name: 'Текст сообщения' })).not.toBeVisible()
  await expect(page.getByRole('textbox', { name: 'Поиск по чатам' })).not.toBeVisible()
}

test('requires only a display name and exposes accessible optional profile fields with no skip action', async ({ page }) => {
  const session = await mockSession(page)
  await page.goto(onboardingPath)

  await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).toBeVisible()
  const displayName = page.getByLabel('Display Name', { exact: true })
  const username = page.getByLabel('Username (optional)', { exact: true })
  const bio = page.getByLabel('Bio (optional)', { exact: true })
  const birthday = page.getByLabel('Birthday (optional)', { exact: true })
  const avatar = page.getByLabel(/Avatar \(optional\)/)
  const submit = page.getByRole('button', { name: 'Create profile', exact: true })

  await expect(displayName).toHaveAttribute('required', '')
  await expect(username).not.toHaveAttribute('required')
  await expect(bio).not.toHaveAttribute('required')
  await expect(birthday).toHaveAttribute('type', 'date')
  await expect(birthday).not.toHaveAttribute('required')
  await expect(avatar).toHaveAttribute('type', 'file')
  await expect(avatar).not.toHaveAttribute('required')
  await expect(avatar).toHaveAttribute('accept', /image\//)
  await expect(page.getByRole('button', { name: /skip|close|dismiss/i })).toHaveCount(0)
  await expect(page.getByRole('link', { name: /skip|home/i })).toHaveCount(0)

  await submit.click()
  await expect(displayName).toBeFocused()
  expect(await displayName.evaluate((input: HTMLInputElement) => input.validity.valueMissing)).toBe(true)
  expect(session.submissions).toBe(0)
  await expectMessengerLocked(page)
})

for (const field of ['Display Name'] as const) {
  test(`rejects a whitespace-only ${field.toLowerCase()} before sending a request`, async ({ page }) => {
    const session = await mockSession(page)
    await page.goto(onboardingPath)
    await fillProfile(page)
    await page.getByLabel(field, { exact: true }).fill('   ')
    await page.getByRole('button', { name: 'Create profile', exact: true }).click()

    await expect(page.getByRole('alert')).toBeVisible()
    await expect(page).toHaveURL(onboardingPath)
    expect(session.submissions).toBe(0)
    await expectMessengerLocked(page)
  })
}

test('creates a display-name-only profile and omits blank optional fields', async ({ page }) => {
  const session = await mockSession(page)
  await page.goto(onboardingPath)
  await page.getByLabel('Display Name', { exact: true }).fill(` ${profile.displayName} `)
  await page.getByLabel('Username (optional)', { exact: true }).fill('   ')
  await page.getByLabel('Bio (optional)', { exact: true }).fill('   ')
  const requestPromise = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()

  expect((await requestPromise).postDataJSON()).toEqual({ displayName: profile.displayName })
  await expect(page).toHaveURL('/')
  expect(session.submissions).toBe(1)
})

test('creates the profile once with trimmed fields and cookies, verifies it after an empty 201 and replaces the route', async ({ page, context, baseURL }) => {
  const session = await mockSession(page)
  await context.addCookies([{
    name: 'refresh_token',
    value: 'onboarding-session',
    url: baseURL!,
    httpOnly: true,
    sameSite: 'Lax',
  }])
  let completeSubmission!: () => void
  const pendingSubmission = new Promise<void>((resolve) => { completeSubmission = resolve })
  session.create = async (route) => {
    await pendingSubmission
    session.hasProfile = true
    await route.fulfill({ status: 201 })
  }

  await page.goto(onboardingPath)
  await fillProfile(page, ` ${profile.displayName} `, ` ${profile.handle} `)
  await page.getByLabel('Bio (optional)', { exact: true }).fill('  Building calm software.\nSay hello.  ')
  await page.getByLabel('Birthday (optional)', { exact: true }).fill('1996-02-29')
  const historyLength = await page.evaluate(() => history.length)
  const readsBeforeCreation = session.profileReads
  const requestPromise = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  const request = await requestPromise

  expect(request.postDataJSON()).toEqual({
    ...profile,
    bio: 'Building calm software.\nSay hello.',
    birthday: '1996-02-29',
  })
  expect(await request.headerValue('content-type')).toBe('application/json')
  expect(await request.headerValue('cookie')).toContain('refresh_token=onboarding-session')
  await expect(page.getByRole('button', { name: 'Creating profile…', exact: true })).toBeDisabled()
  await expect(page.getByLabel('Display Name', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Username (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Bio (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Birthday (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel(/Avatar \(optional\)/)).toBeDisabled()
  await expect(page.getByRole('button', { name: 'Sign out', exact: true })).toBeDisabled()
  await page.keyboard.press('Enter')
  await page.keyboard.press('Enter')
  expect(session.submissions).toBe(1)
  await expect(page).toHaveURL(onboardingPath)
  await expectMessengerLocked(page)

  completeSubmission()
  await expect(page).toHaveURL('/')
  await expect(page.getByRole('button', { name: 'Account menu' })).toBeVisible()
  expect(session.profileReads).toBeGreaterThan(readsBeforeCreation)
  expect(await page.evaluate(() => history.length)).toBe(historyLength)
  await page.getByRole('button', { name: 'Account menu' }).click()
  await expect(page.getByRole('menu').getByText(profile.displayName, { exact: true })).toBeVisible()
  await expect(page.getByRole('menu').getByText(`@${profile.handle}`, { exact: true })).toBeVisible()
  expect(session.submissions).toBe(1)
  expect(session.uploadInitializations).toBe(0)
  expect(session.storageUploads).toBe(0)
  expect(session.uploadCompletions).toBe(0)
})

test('reuses a completed avatar upload after a taken username and accepts a corrected username', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    if (session.submissions === 1) return problem(route, 409, 'HANDLE_ALREADY_EXISTS')
    session.hasProfile = true
    return route.fulfill({ status: 201 })
  }

  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByLabel(/Avatar \(optional\)/).setInputFiles(avatarImage)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText(/username.*taken/i)
  await expect(page.getByLabel('Username (optional)', { exact: true })).toHaveAttribute('aria-invalid', 'true')
  await expect(page.getByLabel('Username (optional)', { exact: true })).toBeFocused()
  await expect(page.getByLabel('Display Name', { exact: true })).toHaveValue(profile.displayName)
  await expect(page.getByLabel('Username (optional)', { exact: true })).toHaveValue(profile.handle)
  await expectMessengerLocked(page)
  expect(session.submissions).toBe(1)

  await page.getByLabel('Username (optional)', { exact: true }).fill('anya_sokolova')
  const requestPromise = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  expect((await requestPromise).postDataJSON()).toEqual({ ...profile, handle: 'anya_sokolova', avatarMediaId: 'media-1' })
  await expect(page).toHaveURL('/')
  expect(session.submissions).toBe(2)
  expect(session.uploadInitializations).toBe(1)
  expect(session.storageUploads).toBe(1)
  expect(session.uploadCompletions).toBe(1)
})

test('retains profile values after a service error and allows a deliberate retry', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    if (session.submissions === 1) return problem(route, 503)
    session.hasProfile = true
    return route.fulfill({ status: 201 })
  }

  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText(/could not create.*profile/i)
  await expect(page.getByRole('alert')).not.toContainText('Internal service detail')
  await expect(page.getByLabel('Display Name', { exact: true })).toHaveValue(profile.displayName)
  await expect(page.getByLabel('Username (optional)', { exact: true })).toHaveValue(profile.handle)
  await expect(page.getByRole('button', { name: 'Create profile', exact: true })).toBeEnabled()
  await expect(page).toHaveURL(onboardingPath)
  expect(session.submissions).toBe(1)

  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page).toHaveURL('/')
  expect(session.submissions).toBe(2)
})

test('previews and removes an avatar, then uploads it through media endpoints before creating the profile', async ({ page, context, baseURL }) => {
  const session = await mockSession(page)
  await context.addCookies([
    { name: 'access_token', value: 'access-token', url: baseURL!, sameSite: 'Lax' },
    { name: 'refresh_token', value: 'refresh-token', url: baseURL!, httpOnly: true, sameSite: 'Lax' },
  ])
  await page.goto(onboardingPath)
  await fillProfile(page)
  const avatar = page.getByLabel(/Avatar \(optional\)/)
  const preview = page.getByRole('img', { name: 'Avatar preview' })

  await avatar.setInputFiles(avatarImage)
  await expect(preview).toBeVisible()
  await expect(preview).toHaveJSProperty('naturalWidth', 1)
  await page.getByRole('button', { name: 'Remove image', exact: true }).click()
  await expect(preview).not.toBeVisible()
  expect(await avatar.evaluate((input: HTMLInputElement) => input.files?.length)).toBe(0)

  await avatar.setInputFiles(avatarImage)
  const initPromise = page.waitForRequest((request) => request.url().endsWith('/api/media/uploads'))
  const putPromise = page.waitForRequest((request) => request.url().includes('/object-storage/upload-1'))
  const completePromise = page.waitForRequest((request) => request.url().endsWith('/api/media/uploads/upload-1/complete'))
  const profilePromise = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  const [initRequest, putRequest, completeRequest, profileRequest] = await Promise.all([
    initPromise, putPromise, completePromise, profilePromise,
  ])

  expect(initRequest.postDataJSON()).toEqual({ fileName: avatarImage.name, purpose: 'avatar' })
  expect(await initRequest.headerValue('authorization')).toBe('Bearer access-token')
  expect(putRequest.method()).toBe('PUT')
  expect(await putRequest.headerValue('content-type')).toBe(avatarImage.mimeType)
  expect(putRequest.postDataBuffer()).toEqual(avatarImage.buffer)
  expect(await putRequest.headerValue('authorization')).toBeNull()
  expect(await putRequest.headerValue('cookie')).toBeNull()
  expect(completeRequest.method()).toBe('POST')
  expect(completeRequest.postData()).toBeNull()
  expect(profileRequest.postDataJSON()).toEqual({ ...profile, avatarMediaId: 'media-1' })
  expect(session.requestOrder).toEqual(['init', 'put', 'complete'])
  await expect(page).toHaveURL('/')
  const savedAvatar = page.getByRole('button', { name: 'Account menu' })
    .getByRole('img', { name: `${profile.displayName} profile photo` })
  await expect(savedAvatar).toHaveJSProperty('naturalWidth', 1)
  await page.reload()
  await expect(savedAvatar).toHaveJSProperty('naturalWidth', 1)
  expect(session.avatarUrlReads).toBe(0)
})

for (const failedStage of ['initialization', 'storage PUT', 'completion'] as const) {
  test(`keeps the avatar and restarts an incomplete upload after ${failedStage} fails`, async ({ page }) => {
    const session = await mockSession(page)
    if (failedStage === 'initialization') {
      session.initUpload = (route) => session.uploadInitializations === 1
        ? problem(route, 503)
        : route.fulfill({ status: 201, json: { id: 'upload-2', presignedUrl: `${storageOrigin}/object-storage/upload-2` } })
    } else if (failedStage === 'storage PUT') {
      session.putUpload = (route) => session.storageUploads === 1
        ? route.fulfill({ status: 503 })
        : route.fulfill({ status: 200 })
    } else {
      session.completeUpload = (route) => session.uploadCompletions === 1
        ? problem(route, 503)
        : route.fulfill({ json: { mediaId: 'media-2' } })
    }

    await page.goto(onboardingPath)
    await fillProfile(page)
    await page.getByLabel(/Avatar \(optional\)/).setInputFiles(avatarImage)
    await page.getByRole('button', { name: 'Create profile', exact: true }).click()

    await expect(page.getByRole('alert')).toHaveText('Could not upload your avatar. Please try again.')
    await expect(page.getByRole('img', { name: 'Avatar preview' })).toBeVisible()
    await expect(page.getByLabel(/Avatar \(optional\)/)).toHaveAttribute('aria-invalid', 'false')
    expect(session.submissions).toBe(0)

    await page.getByRole('button', { name: 'Create profile', exact: true }).click()
    await expect(page).toHaveURL('/')
    expect(session.uploadInitializations).toBe(2)
    expect(session.submissions).toBe(1)
  })
}

for (const invalidAvatar of [
  {
    label: 'an unsupported file type',
    file: { name: 'avatar.svg', mimeType: 'image/svg+xml', buffer: Buffer.from('<svg xmlns="http://www.w3.org/2000/svg"></svg>') },
    error: 'Choose a PNG, JPEG, or WebP image.',
  },
  {
    label: 'an image larger than 2 MB',
    file: { name: 'large-avatar.png', mimeType: 'image/png', buffer: Buffer.alloc(2 * 1024 * 1024 + 1) },
    error: 'Choose an image smaller than 2 MB.',
  },
]) {
  test(`rejects ${invalidAvatar.label} and accepts a valid replacement`, async ({ page }) => {
    const session = await mockSession(page)
    await page.goto(onboardingPath)
    await fillProfile(page)
    await page.getByLabel(/Avatar \(optional\)/).setInputFiles(invalidAvatar.file)

    await expect(page.getByRole('alert')).toContainText(invalidAvatar.error)
    await page.getByLabel('Username (optional)', { exact: true }).press('Enter')
    await expect(page).toHaveURL(onboardingPath)
    expect(session.submissions).toBe(0)
    await expectMessengerLocked(page)

    await page.getByLabel(/Avatar \(optional\)/).setInputFiles(avatarImage)
    await expect(page.getByRole('alert')).not.toBeVisible()
    await page.getByRole('button', { name: 'Create profile', exact: true }).click()
    await expect(page).toHaveURL('/')
    expect(session.submissions).toBe(1)
  })
}

test('a failed read after creation stays locked and retries loading without creating the profile twice', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    session.hasProfile = true
    session.profileStatus = 503
    return route.fulfill({ status: 201 })
  }
  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()

  await expect(page.getByRole('alert')).toContainText('Your profile was created, but could not be loaded. Please try again.')
  await expect(page).toHaveURL(onboardingPath)
  await expectMessengerLocked(page)
  await expect(page.getByLabel('Display Name', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Username (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Bio (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel('Birthday (optional)', { exact: true })).toBeDisabled()
  await expect(page.getByLabel(/Avatar \(optional\)/)).toBeDisabled()
  expect(session.submissions).toBe(1)

  session.profileStatus = undefined
  const readsBeforeRetry = session.profileReads
  await page.getByRole('button', { name: 'Continue to Messenger', exact: true }).click()
  await expect(page).toHaveURL('/')
  await expect(page.getByRole('button', { name: 'Account menu' })).toBeVisible()
  expect(session.profileReads).toBeGreaterThan(readsBeforeRetry)
  expect(session.submissions).toBe(1)
})

test('removing an avatar after profile creation fails omits the completed media ID on retry', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    if (session.submissions === 1) return problem(route, 503)
    session.savedProfile = savedProfileFrom(route)
    session.hasProfile = true
    return route.fulfill({ status: 201 })
  }
  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByLabel(/Avatar \(optional\)/).setInputFiles(avatarImage)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText('Could not create your profile')
  await expectMessengerLocked(page)
  await page.getByRole('button', { name: 'Remove image', exact: true }).click()
  const requestPromise = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  expect((await requestPromise).postDataJSON()).toEqual(profile)
  await expect(page).toHaveURL('/')
  expect(session.submissions).toBe(2)
  expect(session.uploadInitializations).toBe(1)
})

test('replacing an avatar after profile creation fails uploads the replacement', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    if (session.submissions === 1) return problem(route, 503)
    session.savedProfile = savedProfileFrom(route)
    session.hasProfile = true
    return route.fulfill({ status: 201 })
  }
  const replacement = {
    name: 'replacement.jpg',
    mimeType: 'image/jpeg',
    buffer: Buffer.from([255, 216, 255, 217]),
  }
  await page.goto(onboardingPath)
  await fillProfile(page)
  const avatar = page.getByLabel(/Avatar \(optional\)/)
  await avatar.setInputFiles(avatarImage)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page.getByRole('alert')).toContainText('Could not create your profile')

  await avatar.setInputFiles(replacement)
  const profileRequest = page.waitForRequest((request) => request.url().endsWith('/api/me/profile') && request.method() === 'POST')
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  expect((await profileRequest).postDataJSON()).toEqual({ ...profile, avatarMediaId: 'media-2' })
  await expect(page).toHaveURL('/')
  expect(session.uploadInitializations).toBe(2)
})

test('an expired session during creation returns to sign-in without unlocking the messenger', async ({ page }) => {
  const session = await mockSession(page)
  session.create = (route) => {
    session.signedIn = false
    return problem(route, 401)
  }
  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()
  await expect(page).toHaveURL('/auth/login')
  await expect(page.getByRole('heading', { name: 'Welcome back.', level: 1 })).toBeVisible()
  await expectMessengerLocked(page)
  expect(session.submissions).toBe(1)
})

test('an expired session while initializing an avatar upload returns to sign-in', async ({ page }) => {
  const session = await mockSession(page)
  session.initUpload = (route) => {
    session.signedIn = false
    return problem(route, 401)
  }
  await page.goto(onboardingPath)
  await fillProfile(page)
  await page.getByLabel(/Avatar \(optional\)/).setInputFiles(avatarImage)
  await page.getByRole('button', { name: 'Create profile', exact: true }).click()

  await expect(page).toHaveURL('/auth/login')
  await expect(page.getByRole('heading', { name: 'Welcome back.', level: 1 })).toBeVisible()
  await expectMessengerLocked(page)
  expect(session.uploadInitializations).toBe(1)
  expect(session.storageUploads).toBe(0)
  expect(session.submissions).toBe(0)
})

for (const path of ['/', '/auth/login', '/auth/register']) {
  test(`an incomplete account visiting ${path} must finish onboarding, including after a reload`, async ({ page }) => {
    await mockSession(page)
    await page.goto(path)
    await expect(page).toHaveURL(onboardingPath)
    await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).toBeVisible()
    await expectMessengerLocked(page)
    await page.keyboard.press('Escape')
    await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).toBeVisible()

    await page.reload()
    await expect(page).toHaveURL(onboardingPath)
    await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).toBeVisible()
    await expectMessengerLocked(page)
  })
}

test('an account with a profile cannot return to profile creation', async ({ page }) => {
  const session = await mockSession(page, { hasProfile: true })
  await page.goto(onboardingPath)
  await expect(page).toHaveURL('/')
  await expect(page.getByRole('button', { name: 'Account menu' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).not.toBeVisible()
  expect(session.submissions).toBe(0)
})

for (const path of ['/', onboardingPath]) {
  test(`an unauthenticated visit to ${path} requires sign-in`, async ({ page }) => {
    const session = await mockSession(page, { signedIn: false })
    await page.goto(path)
    await expect(page).toHaveURL('/auth/login')
    await expect(page.getByRole('heading', { name: 'Welcome back.', level: 1 })).toBeVisible()
    await expectMessengerLocked(page)
    expect(session.profileReads).toBe(0)
  })
}

for (const status of [403, 503]) {
  test(`a ${status} profile lookup cannot expose the messenger or start a duplicate profile`, async ({ page }) => {
    const session = await mockSession(page, { profileStatus: status })
    await page.goto('/')
    await expect(page.getByRole('button', { name: 'Retry', exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Create your profile.', level: 1 })).not.toBeVisible()
    await expectMessengerLocked(page)
    expect(session.submissions).toBe(0)
  })
}

test('an incomplete account can sign out and cannot access onboarding after ending the session', async ({ page }) => {
  const session = await mockSession(page)
  let logoutRequests = 0
  await page.route('**/api/auth/users/logout', (route) => {
    logoutRequests += 1
    session.signedIn = false
    return route.fulfill({ status: 204 })
  })
  await page.goto(onboardingPath)
  await page.getByRole('button', { name: 'Sign out', exact: true }).click()
  await expect(page).toHaveURL('/auth/login')
  await expect(page.getByRole('heading', { name: 'Welcome back.', level: 1 })).toBeVisible()
  expect(logoutRequests).toBe(1)

  await page.goto(onboardingPath)
  await expect(page).toHaveURL('/auth/login')
  await expectMessengerLocked(page)
})
