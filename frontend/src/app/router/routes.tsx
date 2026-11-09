import type { RouteObject } from 'react-router'
import RegisterPage from '@/pages/auth/Register'
import LoginPage from '@/pages/auth/Login'
import ChatPage from '@/pages/chats/Chat'
import { PageLoading, RootErrorBoundary } from '@/pages/error'
import { ProfileCreationPage } from '@/pages/profile-creation'
import {
  accountLoader,
  onboardingMiddleware,
  protectedMiddleware,
  publicAuthMiddleware,
} from './middleware'
import { ProtectedRouteBoundary } from './ProtectedRouteBoundary'

export const APP_ROUTES: RouteObject[] = [
  {
    path: '/',
    ErrorBoundary: RootErrorBoundary,
    HydrateFallback: PageLoading,
    children: [
      {
        middleware: [protectedMiddleware],
        Component: ProtectedRouteBoundary,
        children: [{ index: true, Component: ChatPage }],
      },
      {
        path: 'auth',
        middleware: [publicAuthMiddleware],
        children: [
          { path: 'register', Component: RegisterPage },
          { path: 'login', Component: LoginPage },
        ],
      },
      {
        path: 'onboarding/profile',
        middleware: [onboardingMiddleware],
        loader: accountLoader,
        Component: ProfileCreationPage,
      },
    ],
  },
]
