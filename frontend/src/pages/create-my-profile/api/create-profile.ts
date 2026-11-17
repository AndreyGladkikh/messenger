import { apiRequest } from '@/shared/api'

interface CreateProfileInput {
  displayName: string
  handle: string
  avatarMediaId?: string
}

export function createProfile({ displayName, handle, avatarMediaId }: CreateProfileInput) {
  return apiRequest<void>('/me/profile', {
    method: 'POST',
    body: {
      displayName,
      handle,
      ...(avatarMediaId ? { avatarMediaId } : {}),
    },
  })
}
