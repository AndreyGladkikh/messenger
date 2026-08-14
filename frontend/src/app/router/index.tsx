import { createBrowserRouter, isRouteErrorResponse, RouterContextProvider, useRouteError } from "react-router";
import RegisterPage from "@/pages/auth/Register";
import ErrorPage from "@/pages/error/GenericError";
import RouteErrorPage from "@/pages/error/RouteError";
import ChatPage from "@/pages/chats/Chat";
import LoginPage from "@/pages/auth/Login";
import { queryClientContext } from "./context";
import queryClient from "@/shared/lib/query-client";
import { authLoader } from "./loaders";

export default createBrowserRouter([
    {
        path: "/",
        ErrorBoundary: RootErrorBoundary,
        children: [
            {
                path: "auth",
                children: [
                    {
                        path: "register",
                        Component: RegisterPage,
                    },
                    {
                        path: "login",
                        Component: LoginPage,
                    },
                ],
            },
            {
                loader: authLoader,
                children: [
                    {
                        index: true,
                        Component: ChatPage,
                    },
                ],
            },
        ]
    }],
    {
        getContext() {
            const context = new RouterContextProvider();
            context.set(queryClientContext, queryClient);
            return context;
        },
    }
);

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