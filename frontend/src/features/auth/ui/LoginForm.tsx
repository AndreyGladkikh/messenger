import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { login } from '../api'
import type { Credentials } from '../model/credentials'
import { CredentialsForm } from './CredentialsForm'

export function LoginForm() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const mutation = useMutation({ mutationFn: login })

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
    navigate('/')
  }

  return (
    <CredentialsForm
      mode="login"
      pending={mutation.isPending}
      error={mutation.error?.message}
      onSubmit={submit}
    />
  )
}
