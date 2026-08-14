import config from "../config";
import http from "../constants/http";

const UNEXPECTED_ERROR_MESSAGE = 'Произошла непредвиденная ошибка.'

export class ApiError extends Error {
    status: number;
    code?: string;
    details?: object;

    constructor(
        status: number,
        message: string,
        code?: string,
        details?: object,
    ) {
        super(message);
        this.name = "ApiError";
        this.status = status;
        this.code = code;
        this.details = details;
    }
}

export class UnauthenticatedError extends ApiError {
    constructor(
        message: string,
        code?: string,
    ) {
        super(
            http.status.unauthorized,
            message || 'Ошибка авторизации.',
            code,
        );
        this.name = "UnauthenticatedError";
    }
}

export class UnexpectedError extends Error {
    constructor(
    ) {
        super(UNEXPECTED_ERROR_MESSAGE);
        this.name = "UnexpectedError";
    }
}

type RequestOptions = Omit<RequestInit, 'body'> & {
    body?: any,
}

export async function apiRequest(path: string, options: RequestOptions = {}, isRefreshTokenIfUnauthorized: boolean = true): Promise<any> {
    const request = await buildRequest(path, options)
    const response = await fetch(request)

    if (response.ok) {
        return await parseSuccess(response)
    }

    if (response.status === http.status.unauthorized && isRefreshTokenIfUnauthorized) {
        try {
            await refreshOnce()
        } catch (e) {
            window.location.href = '/auth/login'
            throw e
        }
        return await apiRequest(path, options, false)
    }

    await parseError(response)
}

async function buildRequest(path: string, options: RequestOptions): Promise<Request> {
    const headers = new Headers({
        [http.header.contentType]: http.contentType.json,
        ...options.headers,
    })

    const accessTokenCookie = await cookieStore.get('access_token');
    if (accessTokenCookie) {
        headers.append('Authorization', `Bearer ${accessTokenCookie.value}`)
    }

    let body = options.body;
    if (body && headers.get(http.header.contentType) === http.contentType.json) {
        body = JSON.stringify(body);
    }

    return new Request(
        `${config.apiUrl}${path}`,
        {
            mode: "cors",
            credentials: "include",
            ...options,
            headers,
            body,
        },
    )
}

async function parseSuccess(response: Response): Promise<any> {
    if (response.status === 204) {
        return null
    }

    const contentType = response.headers.get(http.header.contentType);
    if (contentType?.includes(http.contentType.json)) {
        try {
            const responseData = await response.json()
            return responseData.data
        } catch (e) {
            throw new UnexpectedError()
        }
    }

    try {
        return await response.text()
    } catch (e) {
        throw new UnexpectedError()
    }
}

async function parseError(response: Response): Promise<any> {
    const contentType = response.headers.get(http.header.contentType);
    if (!contentType?.includes(http.contentType.json)) {
        throw new UnexpectedError()
    }
    
    try {
        const responseData = await response.json()

        if (response.status === http.status.unauthorized) {
            throw new UnauthenticatedError(
                responseData?.error?.message,
                responseData?.error?.code,
            )
        }

        throw new ApiError(
            response.status,
            responseData?.error?.message || UNEXPECTED_ERROR_MESSAGE,
            responseData?.error?.code,
            responseData?.error?.details,
        );
    } catch (e) {
        throw new UnexpectedError()
    }
}

let refreshPromise: Promise<void> | null = null

function refreshOnce(): Promise<void> {
    if (!refreshPromise) {
        refreshPromise = refresh().finally(() => { refreshPromise = null })
    }
    return refreshPromise
}

async function refresh(): Promise<void> {
    const response = await fetch(`${config.apiUrl}/auth/refresh`, { method: http.method.post })
    if (!response.ok) {
        await parseError(response)
    }
}