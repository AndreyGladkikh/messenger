import type { ErrorResponse } from "react-router";

export default function RouteErrorPage(error: ErrorResponse) {
    return (
        <>
            <h1>
                {error.status} {error.statusText}
            </h1>
            <p>{error.data}</p>
        </>
    );
}