import { createBrowserRouter } from 'react-router'
import { createRouterContext } from './context'
import { routes } from './routes'

const router = createBrowserRouter(routes, {
  getContext: createRouterContext,
})

export default router
