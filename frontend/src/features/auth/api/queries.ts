import { queryOptions } from "@tanstack/react-query";
import { getCurrentUser } from ".";

export const currentUserQueryKey = ['auth', 'current-user'] as const

export const currentUserQuery = () => queryOptions({
    queryKey: currentUserQueryKey,
    queryFn: getCurrentUser,
})
