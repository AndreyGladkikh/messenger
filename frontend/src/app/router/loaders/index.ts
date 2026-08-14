import { redirect, type LoaderFunctionArgs } from "react-router";
import { queryClientContext } from "../context";
import { currentUserQuery } from "@/features/auth/api/queries";
import { UnauthenticatedError } from "@/shared/api/client";

export async function authLoader({ context }: LoaderFunctionArgs) {
    try {
        const queryClient = context.get(queryClientContext)
        return await queryClient.ensureQueryData(currentUserQuery())
    } catch (e) {
        if (e instanceof UnauthenticatedError) {
            redirect('/auth/login')
            return
        }
        throw e
    }
}