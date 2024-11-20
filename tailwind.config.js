/** @type {import('tailwindcss').Config} */
export default {
    content: ["./resources/**/*.{js,ts,jsx,tsx,html}", "./node_modules/flyonui/dist/js/*.js"],
    theme: {
        extend: {},
    },
    plugins: [require("flyonui"), require("flyonui/plugin")],
};
