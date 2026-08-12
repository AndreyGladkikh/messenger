import { createBrowserRouter, isRouteErrorResponse, useRouteError } from "react-router";
import SignupPage from "@/pages/auth/Signup";
import ErrorPage from "@/pages/error/GenericError";
import RouteErrorPage from "@/pages/error/RouteError";
import ChatPage from "@/pages/chats/Chat";
import SigninPage from "@/pages/auth/Signin";

export default createBrowserRouter([
    {
        path: "/",
        ErrorBoundary: RootErrorBoundary,
        children: [
            {
                index: true,
                Component: ChatPage,
            },
            {
                path: "/signup",
                Component: SignupPage,
            },
            {
                path: "/signin",
                Component: SigninPage,
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