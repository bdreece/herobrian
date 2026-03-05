import { VueQueryPlugin } from '@tanstack/vue-query';

export function install({ app }: Herobrian.Module.Context) {
    app.use(VueQueryPlugin);
}
