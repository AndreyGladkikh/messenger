import { useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { Link } from 'react-router'
import type { Credentials } from '../model/credentials'
import { Card } from '@/shared/ui/pouf/surface'
import { Stack, Row } from '@/shared/ui/pouf/layout'
import { Heading, Text } from '@/shared/ui/pouf/text'
import { Field, Input } from '@/shared/ui/pouf/Input'
import { Button } from '@/shared/ui/pouf/Button'
import { Separator } from '@/shared/ui/pouf/separator'
import { Blob } from '@/shared/ui/pouf/media'
import { ErrorNote } from '@/shared/ui/pouf/feedback'

const DEFAULT_VALUES: Credentials = { login: '', password: '' }

interface CredentialsFormProps {
  mode: 'login' | 'register'
  pending: boolean
  error?: string
  onSubmit: (credentials: Credentials) => Promise<void>
}

export function CredentialsForm({ mode, pending, error, onSubmit }: CredentialsFormProps) {
  const [showPassword, setShowPassword] = useState(false)
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<Credentials>({
    defaultValues: DEFAULT_VALUES,
    mode: 'onSubmit',
    reValidateMode: 'onChange',
    shouldFocusError: true,
  })

  const isRegister = mode === 'register'

  return (
    <main className="grid min-h-[70vh] place-items-center p-6">
      <div className="w-full max-w-[420px]">
        <form onSubmit={handleSubmit(onSubmit)} noValidate>
          <Card>
            <Stack gap={5}>
              <Stack gap={3}>
                <Blob icon="lock" tone="purple" />
                <Stack gap={1}>
                  <Heading level={2}>{isRegister ? 'Create your account' : 'Welcome back'}</Heading>
                  <Text size="sm" muted>
                    {isRegister ? 'Sign up to start using Messenger.' : 'Sign in to your account.'}
                  </Text>
                </Stack>
              </Stack>

              {error ? <ErrorNote>{error}</ErrorNote> : null}

              <Stack gap={4}>
                <Field label="Login" error={errors.login?.message}>
                  {(id, describedBy) => (
                    <Controller
                      name="login"
                      control={control}
                      rules={{ validate: value => value.trim() ? true : 'Login is required.' }}
                      render={({ field }) => (
                        <Input
                          ref={field.ref}
                          id={id}
                          name={field.name}
                          describedBy={describedBy}
                          value={field.value}
                          onChange={field.onChange}
                          onBlur={field.onBlur}
                          placeholder="your-login"
                          autoComplete="username"
                          autoCapitalize="none"
                          spellCheck={false}
                          required
                          invalid={!!errors.login}
                          disabled={pending}
                        />
                      )}
                    />
                  )}
                </Field>

                <Field label="Password" error={errors.password?.message}>
                  {(id, describedBy) => (
                    <Row gap={2} wrap={false}>
                      <div className="min-w-0 flex-1">
                        <Controller
                          name="password"
                          control={control}
                          rules={{ required: 'Password is required.' }}
                          render={({ field }) => (
                            <Input
                              ref={field.ref}
                              id={id}
                              name={field.name}
                              describedBy={describedBy}
                              type={showPassword ? 'text' : 'password'}
                              value={field.value}
                              onChange={field.onChange}
                              onBlur={field.onBlur}
                              placeholder="Password"
                              autoComplete={isRegister ? 'new-password' : 'current-password'}
                              required
                              invalid={!!errors.password}
                              disabled={pending}
                            />
                          )}
                        />
                      </div>
                      <Button
                        variant="quiet"
                        onClick={() => setShowPassword(value => !value)}
                        disabled={pending}
                      >
                        {showPassword ? 'Hide' : 'Show'}
                      </Button>
                    </Row>
                  )}
                </Field>

                <Button block type="submit" loading={pending}>
                  {isRegister ? 'Sign up' : 'Sign in'}
                </Button>
              </Stack>

              <Separator />
              <Row justify="center">
                <Text size="sm" muted>
                  {isRegister ? 'Already have an account? ' : 'New here? '}
                  <Link to={isRegister ? '/auth/login' : '/auth/register'}>
                    {isRegister ? 'Sign in' : 'Create an account'}
                  </Link>
                </Text>
              </Row>
            </Stack>
          </Card>
        </form>
      </div>
    </main>
  )
}
