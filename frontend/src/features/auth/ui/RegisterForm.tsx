import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { register } from '../api'
import type { Credentials } from '../model/credentials'
import { CredentialsForm } from './CredentialsForm'

export function RegisterForm() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const mutation = useMutation({ mutationFn: register })

  async function submit(credentials: Credentials) {
    try {
      await mutation.mutateAsync({
        login: credentials.login.trim(),
        password: credentials.password,
      })
    } catch {
      return
    }
    await queryClient.cancelQueries()
    queryClient.removeQueries()
    navigate('/onboarding/profile')
  }

  return (
    <CredentialsForm
      mode="register"
      pending={mutation.isPending}
      error={mutation.error?.message}
      onSubmit={submit}
    />
  )
}
