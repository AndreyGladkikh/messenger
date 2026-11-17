import { queryOptions } from '@tanstack/react-query'
import { apiRequest, ApiError } from '@/shared/api'
import http from '@/shared/constants/http'
import type { Profile } from '../model/types'

export const profileQueryKey = (accountId: string) => ['profile', accountId] as const
export const profileAvatarQueryKey = (mediaId?: string) => ['profile-avatar', mediaId ?? null] as const

interface DownloadMediaUrl {
  url: string
}

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

export const profileAvatarQuery = (mediaId?: string) => queryOptions({
  queryKey: profileAvatarQueryKey(mediaId),
  queryFn: async () => {
    if (!mediaId) throw new Error('Avatar media ID is required')
    const media = await apiRequest<DownloadMediaUrl>(`/media/${encodeURIComponent(mediaId)}/url`)
    return media.url
  },
  enabled: Boolean(mediaId),
  staleTime: 0,
  gcTime: 0,
  retry: false,
  refetchOnMount: 'always',
  refetchOnWindowFocus: true,
  refetchOnReconnect: true,
})
