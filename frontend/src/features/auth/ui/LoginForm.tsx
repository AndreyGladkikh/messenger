import type { Credentials } from "@/features/auth/model/credentials";
import { CredentialsForm } from "./CredentialsForm";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { login } from "../api";
import { ApiError } from "@/shared/api/client";

export function LoginForm() {
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['me'] })
    },
  })

  const onSubmit = (creds: Credentials) => {
    mutation.mutate(creds)
  }

  const err = mutation.isError ? mutation.error : null
  const isApiErr = err instanceof ApiError

  return (
    <>
      {err && (
        <>
          <div>{err.message}</div>
          {isApiErr && (
            <>
              <div>{err.code}</div>
              <div>{JSON.stringify(err.details, null, 2)}</div>
            </>
          )}
        </>
      )}

      <div>{JSON.stringify(mutation.data, null, 2)}</div>
      <CredentialsForm onSubmit={onSubmit} />
    </>
  )
}
