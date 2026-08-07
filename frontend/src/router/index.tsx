import { createBrowserRouter } from "react-router";
import App from "@/App";
import Test from "@/Test"
import SignupPage from "@/app/signup/page";

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
        path: "/signup",
        Component: SignupPage,
    },
    {
        path: "/signin",
        Component: SignupPage,
    },
  ]);