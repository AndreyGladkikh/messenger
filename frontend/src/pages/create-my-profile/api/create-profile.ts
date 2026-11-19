import { apiRequest } from '@/shared/api'

interface CreateProfileInput {
  displayName: string
  handle?: string
  bio?: string
  birthday?: string
  avatarMediaId?: string
}

export function createProfile({ displayName, handle, bio, birthday, avatarMediaId }: CreateProfileInput) {
  return apiRequest<void>('/me/profile', {
    method: 'POST',
    body: {
      displayName,
      ...(handle ? { handle } : {}),
      ...(bio ? { bio } : {}),
      ...(birthday ? { birthday } : {}),
      ...(avatarMediaId ? { avatarMediaId } : {}),
    },
  })
}
