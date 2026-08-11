const API_URL = import.meta.env.VITE_API_URL;
const JSON_CONTENT_TYPE = "application/json";

export class ApiError extends Error {
    status: number;
    code: string | undefined;
    details: object | undefined;

    constructor(
        status: number,
        code: string | undefined,
        message: string,
        details: object | undefined,
    ) {
        super(message);
        this.name = "ApiError";
        this.status = status;
        this.code = code;
        this.details = details;
    }
}

export async function apiRequest(path: string, options: RequestInit = {}): Promise<any> {
    const headers = new Headers({
        "Content-Type": JSON_CONTENT_TYPE,
        ...options.headers,
    })

    let body = options.body;
    if (body && headers.get("Content-Type") === JSON_CONTENT_TYPE) {
        body = JSON.stringify(body);
    }

    const request = new Request(
        `${API_URL}${path}`,
        {
            mode: "cors",
            credentials: "include",
            ...options,
            headers,
            body,
        },
    )

    const response = await fetch(request)

    if (!response.ok) {
        const responseData = await response.json().catch(() => ({}));
        throw new ApiError(
            response.status,
            responseData.error.code,
            responseData.error.message || "An error occurred",
            responseData.error.details,
        );
    }

    if (response.status === 204) {
        return undefined;
    }

    const responseData = await response.json();
    return new Promise(() => responseData.data);
}