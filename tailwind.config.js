/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./pkg/web/views/**/*.{html,js}",
        "./public/ts/**/*.{ts,js,vue}",
        "./public/pages/**/*.html",
    ],
    plugins: [
        require('@tailwindcss/typography'),
        require("@tailwindcss/forms"),
        require('@tailwindcss/container-queries'),
    ],
}
