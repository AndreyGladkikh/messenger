import type { Credentials } from "@/features/auth/model/credentials";
import { apiRequest } from "@/shared/api/client";

export const login = async (creds: Credentials) => await apiRequest('/auth/login', {
    method: 'POST',
    body: {
        login: creds.login,
        password: creds.password,
    },
})

export const getCurrentUser = async () => await apiRequest('auth/me', {
    method: 'GET',
})