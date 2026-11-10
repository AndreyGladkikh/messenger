import { queryOptions } from '@tanstack/react-query'
import { apiRequest, ApiError } from '@/shared/api'
import http from '@/shared/constants/http'
import type { Profile } from '../model/types'

export const profileQueryKey = (accountId: string) => ['profile', accountId] as const

async function getMyProfile(): Promise<Profile | null> {
  try {
    return await apiRequest<Profile>('/me/profile')
  } catch (error) {
    if (error instanceof ApiError && error.status === http.status.notFound) {
      return null
    }
    throw error
  }
}

export const profileQuery = (accountId: string) => queryOptions({
  queryKey: profileQueryKey(accountId),
  queryFn: getMyProfile,
})
