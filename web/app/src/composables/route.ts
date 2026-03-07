import type { RouteLocationNormalizedLoadedGeneric } from 'vue-router';

export type RouteGetter<T> = (route: RouteLocationNormalizedLoadedGeneric) => T;

export type MaybeRouteGetter<T> = T | RouteGetter<T>;

export function fromRoute<T>(
    route: RouteLocationNormalizedLoadedGeneric,
    value: MaybeRouteGetter<T>,
): T {
    if (typeof value !== 'function') {
        return value;
    }

    return (value as RouteGetter<T>)(route);
}
