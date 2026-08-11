import { createBrowserRouter, isRouteErrorResponse, useRouteError } from "react-router";
import App from "@/App";
import SignupPage from "@/pages/auth/Signup";
import ErrorPage from "@/pages/error/GenericError";
import RouteErrorPage from "@/pages/error/RouteError";
import ChatPage from "@/pages/chats/Chat";

export default createBrowserRouter([
    {
        path: "/",
        ErrorBoundary: RootErrorBoundary,
        children: [
            {
                index: true,
                Component: App,
            },
            {
                path: "/signup",
                Component: SignupPage,
            },
            {
                path: "/signin",
                Component: SignupPage,
            },
            {
                path: "/chat",
                Component: ChatPage,
            },
        ]
    },
  ]);

function RootErrorBoundary() {
    let error = useRouteError();
    if (isRouteErrorResponse(error)) {
        return RouteErrorPage(error);
    } else if (error instanceof Error) {
        return ErrorPage(error);
    } else {
        return <h1>Unknown Error</h1>;
    }
  }