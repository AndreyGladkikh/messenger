import type { QueryClient } from "@tanstack/react-query";
import { createContext } from "react-router";
import { RouterContextProvider } from 'react-router'
import queryClient from '@/shared/lib/query-client'

export const queryClientContext = createContext<QueryClient>()

export function createRouterContext() {
    const context = new RouterContextProvider()
    context.set(queryClientContext, queryClient)
    return context
  }