import { createBrowserRouter } from "react-router";
import App from "../App";

export default createBrowserRouter([
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