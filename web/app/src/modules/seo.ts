export function install({ router }: Herobrian.Module.Context) {
    router.afterEach(to => {
        useSeoMeta(to.meta.seo);
    });
}
