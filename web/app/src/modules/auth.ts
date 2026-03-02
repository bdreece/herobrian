import { useAuth } from '../composables/auth';

export function install({ router }: Herobrian.Module.Context) {
    router.beforeEach((to, _, next) => {
        const [authenticated] = useAuth();
        if (!authenticated.value && !to.meta.allowAnonymous) {
            next({ name: 'login' });
        } else if (authenticated.value && to.name === 'login') {
            next({ name: 'home' });
        } else {
            next();
        }
    });
}
