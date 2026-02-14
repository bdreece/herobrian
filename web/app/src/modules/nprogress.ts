import nprogress from 'nprogress';

export function install({ router }: Herobrian.Module.Context) {
    router.beforeEach((to, from) => {
        if (to.path !== from.path) {
            nprogress.start();
        }
    });

    router.afterEach(() => {
        nprogress.done();
    });
}
