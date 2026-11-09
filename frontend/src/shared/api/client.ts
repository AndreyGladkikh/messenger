import config from "@/shared/config"
import http from "@/shared/constants/http"

const UNEXPECTED_ERROR_MESSAGE = 'Произошла непредвиденная ошибка.'

interface ProblemDetails {
  title?: string
  code?: string
  errors?: object
}

export class ApiError extends Error {
  readonly status: number
  readonly code?: string
  readonly details?: object

  constructor(
    status: number,
    message: string,
    code?: string,
    details?: object,
  ) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}

export class UnauthenticatedError extends ApiError {
  constructor(message?: string, code?: string) {
    super(http.status.unauthorized, message || 'Ошибка авторизации.', code)
    this.name = 'UnauthenticatedError'
  }
}

export class UnexpectedError extends Error {
  constructor() {
    super(UNEXPECTED_ERROR_MESSAGE)
    this.name = 'UnexpectedError'
  }
}

type RequestOptions = Omit<RequestInit, 'body'> & { body?: unknown }

export async function apiRequest<T = unknown>(
  path: string,
  options: RequestOptions = {},
  refreshTokenIfUnauthorized = true,
): Promise<T> {
  const response = await fetch(await buildRequest(path, options))

  if (response.ok) {
    return parseSuccess<T>(response)
  }

  if (response.status === http.status.unauthorized && refreshTokenIfUnauthorized) {
    try {
      await refreshOnce()
    } catch (error) {
      if (
        error instanceof UnauthenticatedError
        || (error instanceof ApiError && error.status === http.status.badRequest)
      ) {
        throw await parseError(response)
      }
      throw error
    }
    return apiRequest<T>(path, options, false)
  }

  throw await parseError(response)
}

async function buildRequest(path: string, options: RequestOptions): Promise<Request> {
  const headers = new Headers({
    [http.header.contentType]: http.contentType.json,
    ...options.headers,
  })

  if ('cookieStore' in window && window.cookieStore) {
    const accessTokenCookie = await window.cookieStore.get(http.cookie.accessToken)
    if (accessTokenCookie) {
      headers.set(http.header.authorization, `Bearer ${accessTokenCookie.value}`)
    }
  }

  const body = options.body && headers.get(http.header.contentType) === http.contentType.json
    ? JSON.stringify(options.body)
    : options.body as BodyInit | null | undefined

  return new Request(`${config.apiUrl}${path}`, {
    mode: 'cors',
    credentials: 'include',
    ...options,
    headers,
    body,
  })
}

async function parseSuccess<T>(response: Response): Promise<T> {
  if (response.status === http.status.noContent || response.headers.get('Content-Length') === '0') {
    return undefined as T
  }

  const contentType = response.headers.get(http.header.contentType)
  try {
    if (contentType?.includes(http.contentType.json)) {
      return await response.json() as T
    }
    return await response.text() as T
  } catch {
    throw new UnexpectedError()
  }
}

async function parseError(response: Response): Promise<ApiError | UnexpectedError> {
  const contentType = response.headers.get(http.header.contentType)
  if (!contentType?.includes(http.contentType.problemPlusJson)) {
    return new UnexpectedError()
  }

  let problem: ProblemDetails
  try {
    problem = await response.json() as ProblemDetails
  } catch {
    return new UnexpectedError()
  }

  if (response.status === http.status.unauthorized) {
    return new UnauthenticatedError(problem.title, problem.code)
  }

  return new ApiError(
    response.status,
    problem.title || UNEXPECTED_ERROR_MESSAGE,
    problem.code,
    problem.errors,
  )
}

let refreshPromise: Promise<unknown> | null = null

function refreshOnce(): Promise<unknown> {
  if (!refreshPromise) {
    refreshPromise = apiRequest(
      '/auth/users/refresh-tokens',
      { method: http.method.post },
      false,
    ).finally(() => { refreshPromise = null })
  }
  return refreshPromise
}
