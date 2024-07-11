import nodeResolve from '@rollup/plugin-node-resolve';
import typescript from '@rollup/plugin-typescript';
import terser from '@rollup/plugin-terser';
import postcss from 'rollup-plugin-postcss';
import autoprefixer from 'autoprefixer';
import cssnano from 'cssnano';
import tailwindcss from 'tailwindcss';

/** @type {import('rollup').RollupOptions} */
export default {
    input: 'src/index.ts',
    output: {
        dir: 'dist',
        format: 'es',
        external: [
            'htmx.org',
            'htmx.org/dist/ext/head-support',
            'htmx.org/dist/ext/sse',
        ],
        paths: {
            'htmx.org': 'https://unpkg.com/htmx.org@1.9.12',
            'htmx.org/dist/ext/head-support': 'https://unpkg.com/htmx.org@1.9.12/dist/ext/head-support.js',
            'htmx.org/dist/ext/sse': 'https://unpkg.com/htmx.org@1.9.12/dist/ext/sse.js',
        }
    },
    plugins: [
        nodeResolve(),
        typescript(),
        terser(),
        postcss({
            extract: true,
            plugins: [
                autoprefixer(),
                cssnano(),
                tailwindcss(),
            ]
        })
    ]
}
