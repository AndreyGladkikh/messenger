import type { RouteObject } from 'react-router'
import { PageLoading, RootErrorBoundary } from '@/pages/error'
import {
  accountLoader,
  onboardingMiddleware,
  protectedMiddleware,
  publicAuthMiddleware,
} from './middleware'
import { LoadingRouteBoundary } from './ProtectedRouteBoundary'
import { HomePage } from '@/pages/home'
import { SignUpPage } from '@/pages/sign-up'
import { SignInPage } from '@/pages/sign-in'
import { CreateMyProfilePage } from '@/pages/create-my-profile'

export const routes: RouteObject[] = [
  {
    Component: LoadingRouteBoundary,
    ErrorBoundary: RootErrorBoundary,
    HydrateFallback: PageLoading,
    children: [
      {
        index: true,
        Component: HomePage,
        middleware: [protectedMiddleware],
      },
      {
        path: 'test',
        Component: SignInPage,
      },
      {
        path: 'auth',
        middleware: [publicAuthMiddleware],
        children: [
          { path: 'register', Component: SignUpPage },
          { path: 'login', Component: SignInPage },
        ],
      },
      {
        path: 'onboarding/profile',
        middleware: [onboardingMiddleware],
        loader: accountLoader,
        Component: CreateMyProfilePage,
      },
    ],
  },
]
