import { addIconSelectors } from '@iconify/tailwind';

/** @type {import('tailwindcss').Config} */
export default {
    content: ["../templates/*.gotmpl", "./src/**/*.ts"],
    theme: {
        extend: {
            colors: {
                'text': '#0e0c06',
                'background': '#fdfcfb',
                'primary': '#b89a4e',
                'secondary': '#b1daa1',
                'accent': '#7fcc81',
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
    plugins: [
        addIconSelectors(['tabler']),
    ],
}

