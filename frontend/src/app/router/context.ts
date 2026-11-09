import type { QueryClient } from "@tanstack/react-query";
import { createContext } from "react-router";

export const queryClientContext = createContext<QueryClient>()