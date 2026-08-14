import { queryOptions } from "@tanstack/react-query";
import { getCurrentUser } from ".";

export const currentUserQuery = () => queryOptions({
    queryKey: ['currentUser'],
    queryFn: getCurrentUser,
})