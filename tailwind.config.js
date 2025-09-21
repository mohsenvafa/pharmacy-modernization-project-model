/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/**/*.templ",
    "./internal/**/*.go",
    "./web/**/*.{html,js}"
  ],
  theme: {
    extend: {}
  },
  plugins: [require("daisyui")]
};
