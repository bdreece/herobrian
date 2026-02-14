import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
    publicDir: '../static/',
    build: {
        assetsDir: '',
        copyPublicDir: false,
        lib: {
            entry: 'src/main.ts',
            name: 'herobrian',
            fileName: 'herobrian',
        },
    },
    plugins: [],
    resolve: {
        alias: {
            '~/': 'src/',
        },
    },
});
