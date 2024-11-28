/** @type {import('tailwindcss').Config} */
export default {
    content: ["./resources/**/*.{js,ts,jsx,tsx,html}", "./node_modules/flyonui/dist/js/*.js", "./node_modules/apexcharts/dist/*.js", "./resources/js/apexcharts-helper.js"],
    theme: {
        extend: {},
    },
    flyonui: {
        themes: ["corporate"],
    },
    plugins: [require("flyonui"), require("flyonui/plugin")],
};
