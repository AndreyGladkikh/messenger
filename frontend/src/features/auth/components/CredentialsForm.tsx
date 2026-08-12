import { useState } from 'react'
import { Controller, useForm } from 'react-hook-form'
import { Card } from '@/components/pouf/surface'
import { Stack, Row } from '@/components/pouf/layout'
import { Heading, Text } from '@/components/pouf/text'
import { Field, Input } from '@/components/pouf/Input'
import { Button } from '@/components/pouf/Button'
import { Separator } from '@/components/pouf/separator'
import { Blob } from '@/components/pouf/media'

interface Credentials {
  login: string
  password: string
}

const DEFAULT_VALUES: Credentials = { login: '', password: '' }

export function CredentialsForm({ onSubmit }: { onSubmit: (creds: Credentials) => void }) {
  const [showPassword, setShowPassword] = useState(false)
  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<Credentials>({
    defaultValues: DEFAULT_VALUES,
    mode: 'onSubmit',
    reValidateMode: 'onChange',
    shouldFocusError: true,
  })

  // const submit = handleSubmit((values) => {
  //   console.log(1, values);
  //   onSubmit(values);
  // })

  return (
    <div style={{ display: 'grid', placeItems: 'center', minHeight: '70vh', padding: 24 }}>
      <div style={{ width: '100%', maxWidth: 420 }}>
        <form onSubmit={handleSubmit(onSubmit)} noValidate>
          <Card>
            <Stack gap={5}>
              <Stack gap={3}>
                <Blob icon='lock' tone="purple" />
                <Stack gap={1}>
                  <Heading level={2}>Welcome back</Heading>
                  <Text size="sm" muted>
                    Sign in to your account.
                  </Text>
                </Stack>
              </Stack>

              <Stack gap={4}>
                <Field label="Login" error={errors.login?.message}>
                  {(id, describedBy) => (
                    <Controller
                      name="login"
                      control={control}
                      rules={{
                        required: 'Login is required.',
                      }}
                      render={({ field }) => (
                        <Input
                          ref={field.ref}
                          id={id}
                          name={field.name}
                          describedBy={describedBy}
                          type="text"
                          value={field.value}
                          onChange={field.onChange}
                          onBlur={field.onBlur}
                          placeholder="login"
                          autoCapitalize="none"
                          spellCheck={false}
                          required
                          invalid={!!errors.login}
                        />
                      )}
                    />
                  )}
                </Field>
                <Field label="Password" error={errors.password?.message}>
                  {(id, describedBy) => (
                    <Row gap={2} wrap={false}>
                      <div style={{ flex: 1, minWidth: 0 }}>
                        <Controller
                          name="password"
                          control={control}
                          shouldUnregister
                          rules={{
                            required: 'Password is required.',
                            minLength: {
                              value: 8,
                              message: 'Use at least 8 characters.',
                            },
                          }}
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
                              placeholder="password..."
                              autoComplete="current-password"
                              required
                              invalid={!!errors.password}
                            />
                          )}
                        />
                      </div>
                      <Button variant="quiet" onClick={() => setShowPassword((s) => !s)}>
                        {showPassword ? 'Hide' : 'Show'}
                      </Button>
                    </Row>
                  )}
                </Field>
                <Button block type="submit" loading={isSubmitting}>Sign in</Button>
              </Stack>

              <Separator />
              <Row justify="center">
                <Text size="sm" muted>New here? Create an account.</Text>
              </Row>
            </Stack>
          </Card>
        </form>
      </div>
    </div>
  )
}
