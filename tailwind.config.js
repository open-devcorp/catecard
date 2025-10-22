/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./pkg/infrastructure/web/templates/**/*.{html,js}",
        "./pkg/web/views/**/*.{html,js}",
        "./public/ts/**/*.{ts,js,vue}",
        "./public/pages/**/*.html",
        "./public/src/**/*.{html,js,ts,vue}"
    ],
    plugins: [
        require('@tailwindcss/typography'),
        require("@tailwindcss/forms"),
        require('@tailwindcss/container-queries'),
    ],
}
