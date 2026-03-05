import { createGlobalState } from '@vueuse/core';
import dayjs from 'dayjs';
import { jwtDecode, type JwtPayload } from 'jwt-decode';

export interface JwtClaims extends Required<JwtPayload> {
    given_name: string;
    family_name: string;
    preferred_username: string;
    picture_url: string;
}

export const useToken = createGlobalState(() =>
    useLocalStorage<string>('access_token', null, {
        writeDefaults: false,
    }),
);

export function useAuth() {
    const token = useToken();

    const claims = computed(() =>
        token.value ? jwtDecode<JwtClaims>(token.value) : null,
    );

    const authenticated = computed(
        () => !!claims.value && dayjs(claims.value.exp * 1000).isAfter(dayjs()),
    );

    return makeDestructurable(
        { authenticated, claims, token } as const,
        [authenticated, claims, token] as const,
    );
}
