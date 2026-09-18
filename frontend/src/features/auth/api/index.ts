import type { Credentials } from "@/features/auth/model/credentials";
import { apiRequest } from "@/shared/api/client";
import http from "@/shared/constants/http";

export const register = async (creds: Credentials) => await apiRequest('/auth/register', {
    method: http.method.post,
    body: {
        login: creds.login,
        password: creds.password,
    },
}, false)

export const login = async (creds: Credentials) => await apiRequest('/auth/login', {
    method: http.method.post,
    body: {
        login: creds.login,
        password: creds.password,
    },
}, false)

export const getCurrentUser = async () => await apiRequest('/auth/me', {
    method: 'GET',
})