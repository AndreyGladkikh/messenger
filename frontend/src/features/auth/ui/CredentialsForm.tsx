import { useRef, useState, type SubmitEvent } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { Eye, EyeOff, LoaderCircle } from 'lucide-react'
import { ApiError } from '@/shared/api'
import { Alert, AlertDescription } from '@/shared/ui/alert'
import { Button } from '@/shared/ui/button'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/shared/ui/field'
import { Input } from '@/shared/ui/input'
import { login, register } from '../api'

interface CredentialsFormProps {
  mode: 'sign-in' | 'sign-up'
  onSuccess: () => void | Promise<void>
}

interface FormError {
  message: string
  loginInvalid?: boolean
  passwordInvalid?: boolean
}

export function CredentialsForm({ mode, onSuccess }: CredentialsFormProps) {
  const queryClient = useQueryClient()
  const submitting = useRef(false)
  const [isPending, setIsPending] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState<FormError | null>(null)
  const isSignUp = mode === 'sign-up'

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault()
    if (submitting.current) return

    const form = event.currentTarget
    const data = new FormData(form)
    submitting.current = true
    setIsPending(true)
    setError(null)

    try {
      // Keep credentials in the form/request, outside the shared query cache.
      await (isSignUp ? register : login)({
        login: String(data.get('login') ?? ''),
        password: String(data.get('password') ?? ''),
      })
      await queryClient.cancelQueries()
      queryClient.clear()
      form.reset()
      await onSuccess()
    } catch (cause) {
      if (cause instanceof ApiError && cause.status === 401 && !isSignUp) {
        setError({
          message: 'Incorrect login or password. Please try again.',
          loginInvalid: true,
          passwordInvalid: true,
        })
      } else if (cause instanceof ApiError && cause.status === 409 && isSignUp) {
        setError({ message: 'This login is already taken. Try another one.', loginInvalid: true })
      } else {
        setError({ message: isSignUp
          ? 'Could not create your account. Please try again.'
          : 'Could not sign in. Please try again.' })
      }
    } finally {
      submitting.current = false
      setIsPending(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} aria-label={isSignUp ? 'Sign up' : 'Sign in'} aria-busy={isPending}>
      <FieldGroup>
        <Field data-disabled={isPending} data-invalid={error?.loginInvalid}>
          <FieldLabel htmlFor={`${mode}-login`}>Login</FieldLabel>
          <Input
            id={`${mode}-login`}
            name="login"
            type="text"
            autoComplete="username"
            autoCapitalize="none"
            autoCorrect="off"
            spellCheck={false}
            placeholder={isSignUp ? 'Choose a login' : 'Enter your login'}
            required
            disabled={isPending}
            aria-invalid={error?.loginInvalid || undefined}
            aria-describedby={error?.loginInvalid ? `${mode}-error` : isSignUp ? `${mode}-login-description` : undefined}
            onChange={() => { if (error) setError(null) }}
            className="h-11"
          />
          {isSignUp ? <FieldDescription id={`${mode}-login-description`}>Use this login whenever you sign in.</FieldDescription> : null}
        </Field>
        <Field data-disabled={isPending} data-invalid={error?.passwordInvalid}>
          <div className="flex items-center justify-between gap-2">
            <FieldLabel htmlFor={`${mode}-password`}>Password</FieldLabel>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
              aria-pressed={showPassword}
              aria-controls={`${mode}-password`}
              disabled={isPending}
              onClick={() => setShowPassword((visible) => !visible)}
            >
              {showPassword ? <EyeOff aria-hidden="true" data-icon="inline-start" /> : <Eye aria-hidden="true" data-icon="inline-start" />}
              {showPassword ? 'Hide' : 'Show'}
            </Button>
          </div>
          <Input
            id={`${mode}-password`}
            name="password"
            type={showPassword ? 'text' : 'password'}
            autoComplete={isSignUp ? 'new-password' : 'current-password'}
            placeholder={isSignUp ? 'Create a password' : 'Enter your password'}
            required
            disabled={isPending}
            aria-invalid={error?.passwordInvalid || undefined}
            aria-describedby={error?.passwordInvalid ? `${mode}-error` : undefined}
            onChange={() => { if (error) setError(null) }}
            className="h-11"
          />
        </Field>
        {error ? (
          <Alert variant="destructive" id={`${mode}-error`}>
            <AlertDescription>{error.message}</AlertDescription>
          </Alert>
        ) : null}
        <Field>
          <Button type="submit" size="lg" disabled={isPending} className="h-11 w-full">
            {isPending ? <LoaderCircle aria-hidden="true" data-icon="inline-start" className="animate-spin" /> : null}
            {isPending ? (isSignUp ? 'Creating account…' : 'Signing in…') : (isSignUp ? 'Sign Up' : 'Sign In')}
          </Button>
        </Field>
      </FieldGroup>
    </form>
  )
}
