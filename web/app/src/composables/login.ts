import type { UnwrapRef } from 'vue';

import { useMutation, type UseMutationOptions } from '@tanstack/vue-query';
import axios from 'axios';

export interface LoginResult {
    accessToken: string;
}

export type UseLoginOptions = Omit<
    UnwrapRef<
        Exclude<
            UseMutationOptions<LoginResult, Error, URLSearchParams, unknown>,
            (...args: any[]) => any
        >
    >,
    'mutationFn'
>;

export function useLogin(options: UseLoginOptions) {
    return useMutation({
        ...options,
        mutationFn: login,
    });
}

async function login(params: URLSearchParams) {
    const res = await axios.post<LoginResult>('/api/user/login', params, {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
    });

    if (res.status !== 200) {
        throw new Error('Username or password is incorrect!');
    }

    return res.data;
}
