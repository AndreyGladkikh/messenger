import { createBrowserRouter } from "react-router";
import App from "@/App";
import Test from "@/Test"

export default createBrowserRouter([
    {
        path: "/test",
        element: <div>
            <h1 className="text-3xl font-bold underline">
                Hello world!
            </h1>
            <Test />
        </div>,
    },
    {
        path: "/",
        Component: App,
    },
    {
        path: "/register",
        element: <div>Hello World</ div >,
    },
    {
        path: "/login",
        element: <div>Hello World</ div >,
    },
  ]);