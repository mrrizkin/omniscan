/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./resources/**/*.{js,ts,jsx,tsx,html}",
        "./node_modules/flyonui/dist/js/*.js",
        "./node_modules/apexcharts/dist/*.js",
        "./resources/js/apexcharts-helper.js",
        "./app/providers/app/menu.go",
    ],
    theme: {
        extend: {},
    },
    flyonui: {
        themes: ["light", "dark", "gourmet", "corporate", "luxury", "soft"],
    },
    plugins: [require("flyonui"), require("flyonui/plugin")],
};
