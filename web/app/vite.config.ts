import TailwindCSS from '@tailwindcss/vite';
import { unheadVueComposablesImports } from '@unhead/vue';
import Vue from '@vitejs/plugin-vue';
import AutoImport from 'unplugin-auto-import/vite';
import VueComponents from 'unplugin-vue-components/vite';
import { VueRouterAutoImports } from 'unplugin-vue-router';
import VueRouter from 'unplugin-vue-router/vite';
import { defineConfig } from 'vite';
import VueLayouts from 'vite-plugin-vue-layouts-next';
import VueMacros from 'vue-macros/vite';

// https://vite.dev/config/
export default defineConfig({
    publicDir: '../static/',
    build: {
        assetsDir: '',
        copyPublicDir: true,
    },
    plugins: [
        VueMacros({
            plugins: {
                vue: Vue(),
                vueRouter: VueRouter({
                    routesFolder: 'src/pages',
                    dts: 'src/typed-router.d.ts',
                }),
            },
        }),

        VueLayouts({
            inheritDefaultLayout: false,
            layoutsDirs: ['src/layouts'],
            pagesDirs: ['src/pages'],
        }),

        AutoImport({
            imports: [
                'vue',
                '@vueuse/core',
                unheadVueComposablesImports,
                VueRouterAutoImports,
                {
                    'vue-router/auto': ['useLink'],
                },
            ],
            dirs: ['src/composables'],
            dts: 'src/auto-imports.d.ts',
            vueTemplate: true,
        }),

        VueComponents({
            dirs: ['src/components'],
            dts: 'src/components.d.ts',
        }),

        TailwindCSS(),
    ],
    resolve: {
        alias: {
            '~/': 'src/',
        },
    },
});
