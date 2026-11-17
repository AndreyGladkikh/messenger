export const avatarMaxBytes = 2 * 1024 * 1024
export const avatarAcceptedTypes = ['image/jpeg', 'image/png', 'image/webp']

export interface ProfileFormErrors {
  displayName?: string
  username?: string
  avatar?: string
  form?: string
}

export function validateAvatar(file: File): string | undefined {
  if (!avatarAcceptedTypes.includes(file.type)) {
    return 'Choose a PNG, JPEG, or WebP image.'
  }
  if (!file.size || file.size > avatarMaxBytes) {
    return 'Choose an image smaller than 2 MB.'
  }
}
