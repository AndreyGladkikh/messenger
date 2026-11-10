import { useMutation, useQueryClient } from '@tanstack/react-query'
import { logout } from '../api'

export function useSignOut() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: logout,
    onSuccess: async () => {
      await queryClient.cancelQueries()
      queryClient.clear()
    },
  })
}
