import { createGlobalState } from '@vueuse/core';
import { jwtDecode, type JwtPayload } from 'jwt-decode';

export interface JwtClaims extends Required<JwtPayload> {}

export const useAuth = createGlobalState(() => {
    const token = useLocalStorage<string>('access_token', null, {
        writeDefaults: false,
    });

    const authenticated = computed(() => !!token.value);

    const claims = computed(() =>
        token.value ? jwtDecode<JwtClaims>(token.value) : null,
    );

    return makeDestructurable(
        { authenticated, claims, token } as const,
        [authenticated, claims, token] as const,
    );
});
