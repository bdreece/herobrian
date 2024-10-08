import { addIconSelectors } from '@iconify/tailwind';

/** @type {import('tailwindcss').Config} */
export default {
    content: {
        relative: false,
        files: ['../templates/*.gotmpl', '../../internal/controller/*.go', './src/**/*.ts'],
    },
    theme: {
        extend: {
            colors: {
                text: '#0e0c06',
                background: '#fdfcfb',
                primary: '#b89a4e',
                'primary-dark': '#937B3E',
                secondary: '#b1daa1',
                'secondary-dark': '#8EAE81',
                accent: '#7fcc81',
            },
            fontSize: {
                sm: '0.750rem',
                base: '1rem',
                xl: '1.333rem',
                '2xl': '1.777rem',
                '3xl': '2.369rem',
                '4xl': '3.158rem',
                '5xl': '4.210rem',
            },
            fontFamily: {
                sans: ['Montserrat', 'system-ui', 'sans-serif'],
            },
            fontWeight: {
                normal: '400',
                bold: '700',
            },
        },
    },
    plugins: [addIconSelectors(['tabler'])],
};
