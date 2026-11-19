import { useEffect, useRef, useState, type ChangeEvent, type SubmitEvent } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { ArrowRight, LoaderCircle, UserRound } from 'lucide-react'
import { useLoaderData, useNavigate } from 'react-router'
import { profileQuery } from '@/entities/profile'
import { useSignOut, type Account } from '@/features/auth'
import { ApiError, UnauthenticatedError } from '@/shared/api'
import { Alert, AlertDescription } from '@/shared/ui/alert'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/shared/ui/card'
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/shared/ui/field'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { createProfile } from '../api/create-profile'
import { AvatarUploadError, uploadAvatar } from '../api/upload-avatar'
import { avatarAcceptedTypes, validateAvatar, type ProfileFormErrors } from '../model/profile-form'

export function CreateMyProfilePage() {
  const account = useLoaderData<Account>()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const signOut = useSignOut()
  const submitting = useRef(false)
  const uploadedAvatar = useRef<{ file: File; mediaId: string } | null>(null)
  const focusAfterSubmit = useRef<HTMLInputElement | null>(null)
  const fileInput = useRef<HTMLInputElement>(null)
  const [isPending, setIsPending] = useState(false)
  const [profileCreated, setProfileCreated] = useState(false)
  const [avatar, setAvatar] = useState<{ file: File; preview: string } | null>(null)
  const [errors, setErrors] = useState<ProfileFormErrors>({})
  const busy = isPending || signOut.isPending

  useEffect(() => {
    return () => { if (avatar) URL.revokeObjectURL(avatar.preview) }
  }, [avatar])

  useEffect(() => {
    if (!isPending) {
      focusAfterSubmit.current?.focus()
      focusAfterSubmit.current = null
    }
  }, [isPending])

  function clearError(field: keyof ProfileFormErrors) {
    setErrors((current) => ({ ...current, [field]: undefined, form: undefined }))
  }

  function selectAvatar(event: ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.[0]
    if (!file) return
    const error = validateAvatar(file)
    uploadedAvatar.current = null
    setErrors((current) => ({ ...current, avatar: error, form: undefined }))
    setAvatar(error ? null : { file, preview: URL.createObjectURL(file) })
    if (error) event.currentTarget.value = ''
  }

  async function finishOnboarding() {
    // Creation returns no body. Read the persisted profile before unlocking the app.
    const profile = await queryClient.fetchQuery(profileQuery(account.id))
    if (!profile) throw new Error('Profile is not available yet')
    await navigate('/', { replace: true })
  }

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault()
    if (submitting.current || busy) return

    const form = event.currentTarget
    const data = new FormData(form)
    const displayName = String(data.get('displayName') ?? '').trim()
    const handle = String(data.get('username') ?? '').trim() || undefined
    const bio = String(data.get('bio') ?? '').trim() || undefined
    const birthday = String(data.get('birthday') ?? '') || undefined
    const validation: ProfileFormErrors = {
      displayName: displayName ? undefined : 'Enter your display name.',
      avatar: errors.avatar,
    }
    if (!profileCreated && Object.values(validation).some(Boolean)) {
      setErrors(validation)
      const field = !displayName ? 'displayName' : 'avatar'
      const input = form.elements.namedItem(field) as HTMLInputElement | null
      input?.focus()
      return
    }

    submitting.current = true
    setIsPending(true)
    setErrors({})
    let created = profileCreated

    try {
      if (!created) {
        const completedAvatar = uploadedAvatar.current
        let avatarMediaId = completedAvatar?.file === avatar?.file
          ? completedAvatar?.mediaId
          : undefined
        if (avatar && !avatarMediaId) {
          avatarMediaId = await uploadAvatar(avatar.file)
          uploadedAvatar.current = { file: avatar.file, mediaId: avatarMediaId }
        }
        await createProfile({ displayName, handle, bio, birthday, avatarMediaId })
        created = true
        setProfileCreated(true)
      }
      await finishOnboarding()
    } catch (error) {
      if (error instanceof UnauthenticatedError) {
        await queryClient.cancelQueries()
        queryClient.clear()
        await navigate('/auth/login', { replace: true })
      } else if (created) {
        setErrors({ form: 'Your profile was created, but could not be loaded. Please try again.' })
      } else if (error instanceof AvatarUploadError) {
        setErrors({ form: 'Could not upload your avatar. Please try again.' })
      } else if (error instanceof ApiError && error.code === 'HANDLE_ALREADY_EXISTS') {
        setErrors({ username: 'This username is already taken. Try another one.' })
        focusAfterSubmit.current = form.elements.namedItem('username') as HTMLInputElement | null
      } else {
        setErrors({ form: 'Could not create your profile. Please try again.' })
      }
    } finally {
      submitting.current = false
      setIsPending(false)
    }
  }

  function handleSignOut() {
    if (busy) return
    signOut.mutate(undefined, {
      onSuccess: () => { void navigate('/auth/login', { replace: true }) },
    })
  }

  return (
    <main className="not-typeset flex min-h-svh flex-col items-center justify-center bg-muted/40 px-4 py-10 font-[family-name:var(--font-geist)]">
      <title>Create your profile · Messenger</title>
      <div className="flex w-full max-w-md flex-col gap-7">
        <div className="flex items-center justify-center gap-2.5 font-medium">
          <span aria-hidden="true" className="flex size-9 items-center justify-center rounded-xl bg-primary font-mono text-primary-foreground">m</span>
          Messenger
        </div>
        <Card className="gap-6 pt-6">
          <CardHeader className="gap-2 px-6">
            <CardTitle><h1 className="text-2xl font-semibold tracking-tight">Create your profile.</h1></CardTitle>
            <CardDescription>One last step before your first conversation. Add a name and any details you want to share.</CardDescription>
          </CardHeader>
          <CardContent className="px-6">
            <form onSubmit={handleSubmit} aria-label="Create profile" aria-busy={isPending}>
              <FieldGroup>
                <Field data-disabled={busy || profileCreated} data-invalid={Boolean(errors.displayName)}>
                  <FieldLabel htmlFor="profile-display-name">Display Name</FieldLabel>
                  <Input
                    id="profile-display-name" name="displayName" autoComplete="nickname"
                    placeholder="What should we call you?" required className="h-11"
                    disabled={busy || profileCreated} aria-invalid={Boolean(errors.displayName)}
                    aria-describedby={errors.displayName ? 'profile-display-name-error' : undefined}
                    onChange={() => clearError('displayName')}
                  />
                  {errors.displayName ? <FieldError id="profile-display-name-error">{errors.displayName}</FieldError> : null}
                </Field>
                <Field data-disabled={busy || profileCreated} data-invalid={Boolean(errors.username)}>
                  <FieldLabel htmlFor="profile-username">Username (optional)</FieldLabel>
                  <Input
                    id="profile-username" name="username" autoComplete="username"
                    autoCapitalize="none" autoCorrect="off" spellCheck={false}
                    placeholder="Choose a unique username" className="h-11"
                    disabled={busy || profileCreated} aria-invalid={Boolean(errors.username)}
                    aria-describedby={`profile-username-description${errors.username ? ' profile-username-error' : ''}`}
                    onChange={() => clearError('username')}
                  />
                  <FieldDescription id="profile-username-description">Your public username helps others find you.</FieldDescription>
                  {errors.username ? <FieldError id="profile-username-error">{errors.username}</FieldError> : null}
                </Field>
                <Field data-disabled={busy || profileCreated}>
                  <FieldLabel htmlFor="profile-bio">Bio (optional)</FieldLabel>
                  <Textarea
                    id="profile-bio" name="bio" placeholder="Tell people a little about yourself"
                    disabled={busy || profileCreated} aria-describedby="profile-bio-description"
                  />
                  <FieldDescription id="profile-bio-description">A short introduction shown on your profile.</FieldDescription>
                </Field>
                <Field data-disabled={busy || profileCreated}>
                  <FieldLabel htmlFor="profile-birthday">Birthday (optional)</FieldLabel>
                  <Input
                    id="profile-birthday" name="birthday" type="date" className="h-11"
                    disabled={busy || profileCreated} aria-describedby="profile-birthday-description"
                  />
                  <FieldDescription id="profile-birthday-description">Stored as a calendar date without timezone conversion.</FieldDescription>
                </Field>
                <Field data-disabled={busy || profileCreated} data-invalid={Boolean(errors.avatar)}>
                  <FieldLabel htmlFor="profile-avatar">Avatar (optional)</FieldLabel>
                  <div className="flex items-center gap-4">
                    <div className="flex size-14 shrink-0 items-center justify-center overflow-hidden rounded-full bg-muted text-muted-foreground">
                      {avatar ? <img src={avatar.preview} alt="Avatar preview" className="size-full object-cover" /> : <UserRound aria-hidden="true" className="size-6" />}
                    </div>
                    <div className="flex min-w-0 flex-1 flex-col gap-2">
                      <Input
                        ref={fileInput} id="profile-avatar" name="avatar" type="file"
                        accept={avatarAcceptedTypes.join(',')} onChange={selectAvatar}
                        disabled={busy || profileCreated} aria-invalid={Boolean(errors.avatar)}
                        aria-describedby={`profile-avatar-description${errors.avatar ? ' profile-avatar-error' : ''}`}
                      />
                      {avatar || errors.avatar ? (
                        <Button type="button" variant="ghost" size="sm" className="self-start" disabled={busy || profileCreated}
                          onClick={() => {
                            uploadedAvatar.current = null
                            setAvatar(null)
                            clearError('avatar')
                            if (fileInput.current) fileInput.current.value = ''
                          }}>
                          Remove image
                        </Button>
                      ) : null}
                    </div>
                  </div>
                  <FieldDescription id="profile-avatar-description">PNG, JPEG, or WebP. Up to 2 MB.</FieldDescription>
                  {errors.avatar ? <FieldError id="profile-avatar-error">{errors.avatar}</FieldError> : null}
                </Field>
                {errors.form ? <Alert variant="destructive"><AlertDescription>{errors.form}</AlertDescription></Alert> : null}
                <Field>
                  <Button type="submit" size="lg" disabled={busy} className="h-11 w-full">
                    {isPending ? <LoaderCircle aria-hidden="true" data-icon="inline-start" className="animate-spin" /> : null}
                    {isPending ? 'Creating profile…' : profileCreated ? 'Continue to Messenger' : 'Create profile'}
                    {!isPending ? <ArrowRight aria-hidden="true" data-icon="inline-end" /> : null}
                  </Button>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
          <CardFooter className="flex-col gap-2 px-6 py-4">
            <span className="text-center text-muted-foreground">Complete your profile to use Messenger.</span>
            <Button type="button" variant="ghost" size="sm" disabled={busy} onClick={handleSignOut}>
              {signOut.isPending ? 'Signing out…' : 'Sign out'}
            </Button>
            {signOut.isError ? <Alert variant="destructive"><AlertDescription>Could not sign out. Please try again.</AlertDescription></Alert> : null}
          </CardFooter>
        </Card>
      </div>
    </main>
  )
}
