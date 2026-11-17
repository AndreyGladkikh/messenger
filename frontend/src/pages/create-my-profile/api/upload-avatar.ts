import { apiRequest, UnauthenticatedError } from '@/shared/api'

interface InitedMediaUpload {
  id: string
  presignedUrl: string
}

interface Media {
  mediaId: string
}

export class AvatarUploadError extends Error {
  readonly cause: unknown

  constructor(cause: unknown) {
    super('Could not upload avatar')
    this.name = 'AvatarUploadError'
    this.cause = cause
  }
}

export async function uploadAvatar(file: File): Promise<string> {
  try {
    const upload = await apiRequest<InitedMediaUpload>('/media/uploads', {
      method: 'POST',
      body: { fileName: file.name, purpose: 'avatar' },
    })

    const response = await fetch(upload.presignedUrl, {
      method: 'PUT',
      body: file,
      headers: { 'Content-Type': file.type },
      credentials: 'omit',
    })
    if (!response.ok) {
      throw new Error(`Object storage rejected the upload with status ${response.status}`)
    }

    const media = await apiRequest<Media>(`/media/uploads/${encodeURIComponent(upload.id)}/complete`, {
      method: 'POST',
    })
    return media.mediaId
  } catch (error) {
    if (error instanceof UnauthenticatedError) throw error
    throw new AvatarUploadError(error)
  }
}
