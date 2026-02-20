import { useAuth } from '../composables/auth';

export function install({ router }: Herobrian.Module.Context) {
    router.beforeEach((to, _, next) => {
        const [authenticated] = useAuth();
        if (!authenticated.value && !to.meta.allowAnonymous) {
            next({ name: '/user/(auth)/login' });
        } else if (authenticated.value && to.name === '/user/(auth)/login') {
            next({ name: '/' });
        } else {
            next();
        }
    });
}
