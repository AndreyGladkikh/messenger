import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Controller, useForm, useWatch } from 'react-hook-form'
import { useLoaderData, useNavigate } from 'react-router'
import type { Account } from '@/features/auth'
import { profileQuery, profileQueryKey, type Profile } from '@/entities/profile'
import { ApiError } from '@/shared/api/client'
import http from '@/shared/constants/http'
import { Avatar } from '@/shared/ui/pouf/avatar'
import { Button } from '@/shared/ui/pouf/Button'
import { ErrorNote } from '@/shared/ui/pouf/feedback'
import { Field, Input } from '@/shared/ui/pouf/Input'
import { Stack, Row } from '@/shared/ui/pouf/layout'
import { Blob } from '@/shared/ui/pouf/media'
import { Card } from '@/shared/ui/pouf/surface'
import { Heading, Text } from '@/shared/ui/pouf/text'
import { createProfile } from '../api/create-profile'

interface ProfileFormValues {
  displayName: string
  handle: string
  avatarUrl: string
}

const DEFAULT_VALUES: ProfileFormValues = {
  displayName: '',
  handle: '',
  avatarUrl: '',
}

function validAvatarUrl(value: string): string | undefined {
  if (!value.trim()) return undefined
  try {
    const url = new URL(value.trim())
    return url.protocol === 'http:' || url.protocol === 'https:' ? url.toString() : undefined
  } catch {
    return undefined
  }
}

export function ProfileCreationPage() {
  const account = useLoaderData<Account>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const mutation = useMutation({ mutationFn: createProfile })
  const {
    control,
    handleSubmit,
    setError,
    clearErrors,
    formState: { errors, isSubmitting },
  } = useForm<ProfileFormValues>({
    defaultValues: DEFAULT_VALUES,
    mode: 'onSubmit',
    reValidateMode: 'onChange',
    shouldFocusError: true,
  })
  const displayName = useWatch({ control, name: 'displayName' })
  const avatarUrl = useWatch({ control, name: 'avatarUrl' })
  const avatarPreview = validAvatarUrl(avatarUrl)

  async function finish(profile: Profile) {
    queryClient.setQueryData(profileQueryKey(account.id), profile)
    await navigate('/', { replace: true })
  }

  async function submit(values: ProfileFormValues) {
    clearErrors('root')
    const normalizedAvatarUrl = validAvatarUrl(values.avatarUrl)
    const profile: Profile = {
      displayName: values.displayName.trim(),
      handle: values.handle.trim(),
      ...(normalizedAvatarUrl ? { avatarUrl: normalizedAvatarUrl } : {}),
    }

    try {
      await mutation.mutateAsync(profile)
      await finish(profile)
      return
    } catch (error) {
      if (!(error instanceof ApiError) || error.status !== http.status.conflict) {
        setError('root.server', { message: error instanceof Error ? error.message : 'Could not create your profile.' })
        return
      }

      try {
        queryClient.removeQueries({ queryKey: profileQueryKey(account.id), exact: true })
        const existing = await queryClient.fetchQuery(profileQuery(account.id))
        if (existing) {
          await finish(existing)
          return
        }
      } catch (recheckError) {
        setError('root.server', {
          message: recheckError instanceof Error ? recheckError.message : 'Could not confirm your profile. Please retry.',
        })
        return
      }

      if (error.code === 'HANDLE_ALREADY_EXISTS') {
        setError('handle', { message: 'This username is already taken.' }, { shouldFocus: true })
      } else {
        setError('root.server', { message: 'A profile already exists for this account. Please retry.' })
      }
    }
  }

  return (
    <main className="grid min-h-screen place-items-center p-6">
      <div className="w-full max-w-[520px]">
        <form onSubmit={handleSubmit(submit)} noValidate>
          <Card>
            <Stack gap={5}>
              <Row gap={3} wrap={false}>
                <Blob icon="target" tone="purple" />
                <Stack gap={1}>
                  <Heading level={2}>Create your profile</Heading>
                  <Text size="sm" muted>Your profile is required before you can use Messenger.</Text>
                </Stack>
              </Row>

              <Row justify="center">
                <Avatar
                  src={avatarPreview}
                  alt={avatarPreview ? `${displayName || 'Your'} avatar preview` : undefined}
                  fallback={displayName.trim().slice(0, 2) || 'ME'}
                  size="lg"
                />
              </Row>

              {errors.root?.server?.message ? <ErrorNote>{errors.root.server.message}</ErrorNote> : null}

              <Stack gap={4}>
                <Field label="Display Name" error={errors.displayName?.message}>
                  {(id, describedBy) => (
                    <Controller
                      name="displayName"
                      control={control}
                      rules={{ validate: value => value.trim() ? true : 'Display name is required.' }}
                      render={({ field }) => (
                        <Input
                          ref={field.ref}
                          id={id}
                          name={field.name}
                          describedBy={describedBy}
                          value={field.value}
                          onChange={field.onChange}
                          onBlur={field.onBlur}
                          placeholder="How people will see you"
                          autoComplete="name"
                          required
                          invalid={!!errors.displayName}
                          disabled={isSubmitting}
                        />
                      )}
                    />
                  )}
                </Field>

                <Field label="Username / Handle" error={errors.handle?.message}>
                  {(id, describedBy) => (
                    <Controller
                      name="handle"
                      control={control}
                      rules={{ validate: value => value.trim() ? true : 'Username is required.' }}
                      render={({ field }) => (
                        <Input
                          ref={field.ref}
                          id={id}
                          name={field.name}
                          describedBy={describedBy}
                          value={field.value}
                          onChange={field.onChange}
                          onBlur={field.onBlur}
                          placeholder="your-handle"
                          autoComplete="username"
                          autoCapitalize="none"
                          spellCheck={false}
                          required
                          invalid={!!errors.handle}
                          disabled={isSubmitting}
                        />
                      )}
                    />
                  )}
                </Field>

                <Field
                  label="Avatar URL (optional)"
                  hint="Use an absolute http:// or https:// image URL."
                  error={errors.avatarUrl?.message}
                >
                  {(id, describedBy) => (
                    <Controller
                      name="avatarUrl"
                      control={control}
                      rules={{
                        validate: value => !value.trim() || validAvatarUrl(value)
                          ? true
                          : 'Enter an absolute http:// or https:// URL.',
                      }}
                      render={({ field }) => (
                        <Input
                          ref={field.ref}
                          id={id}
                          name={field.name}
                          describedBy={describedBy}
                          type="url"
                          value={field.value}
                          onChange={field.onChange}
                          onBlur={field.onBlur}
                          placeholder="https://example.com/avatar.jpg"
                          autoComplete="url"
                          inputMode="url"
                          autoCapitalize="none"
                          spellCheck={false}
                          invalid={!!errors.avatarUrl}
                          disabled={isSubmitting}
                        />
                      )}
                    />
                  )}
                </Field>

                <Button block type="submit" loading={isSubmitting}>Create profile</Button>
              </Stack>
            </Stack>
          </Card>
        </form>
      </div>
    </main>
  )
}
