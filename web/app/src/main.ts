import { createHead } from '@unhead/vue/client';
import { setupLayouts } from 'virtual:generated-layouts';
import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import { routes } from 'vue-router/auto-routes';
import App from './App.vue';
import './styles/index.css';

const app = createApp(App);
const head = createHead();
const router = createRouter({
    history: createWebHistory(),
    routes: setupLayouts(routes),
});

app.use(router).use(router);

const modules = import.meta.glob<Herobrian.Module>('./modules/*.ts', {
    eager: true,
});

const context = {
    app,
    head,
    router,
    routes,
};

await Promise.all(
    Object.values(modules).map(module =>
        Promise.resolve(module.install?.(context)),
    ),
);

app.mount('#app');
