import type { Account, Credentials } from "../model/credentials";
import { apiRequest } from "@/shared/api";
import http from "@/shared/constants/http";

export const register = async (creds: Credentials) => await apiRequest('/auth/users', {
    method: http.method.post,
    body: {
        login: creds.login,
        password: creds.password,
    },
}, false)

export const login = async (creds: Credentials) => await apiRequest('/auth/users/login', {
    method: http.method.post,
    body: {
        login: creds.login,
        password: creds.password,
    },
}, false)

export const logout = () => apiRequest<void>('/auth/users/logout', {
    method: http.method.post,
}, false)

export const getCurrentUser = async () => await apiRequest<Account>('/auth/me', {
    method: 'GET',
})
