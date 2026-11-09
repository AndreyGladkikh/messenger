import type { Profile } from '@/entities/profile'
import { apiRequest } from '@/shared/api/client'
import http from '@/shared/constants/http'

export async function createProfile(profile: Profile): Promise<void> {
  await apiRequest<void>('/me/profile', {
    method: http.method.post,
    body: profile,
  })
}
