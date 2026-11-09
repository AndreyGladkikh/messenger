import { createBrowserRouter, RouterContextProvider } from 'react-router'
import queryClient from '@/shared/lib/query-client'
import { queryClientContext } from './context'
import { APP_ROUTES } from './routes'

export function createRouterContext() {
  const context = new RouterContextProvider()
  context.set(queryClientContext, queryClient)
  return context
}

const router = createBrowserRouter(APP_ROUTES, {
  future: { v8_middleware: true },
  getContext: createRouterContext,
})

export default router
