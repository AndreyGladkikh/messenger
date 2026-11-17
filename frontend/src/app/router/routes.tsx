import type { RouteObject } from 'react-router'
import { PageLoading, RootErrorBoundary } from '@/pages/error'
import {
  accountLoader,
  homeProfileLoader,
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
        middleware: [protectedMiddleware],
        children: [
          { index: true, loader: homeProfileLoader, Component: HomePage },
        ],
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
