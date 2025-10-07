const path = require("path");

const fromWeb = (...segments) => path.join(__dirname, ...segments);

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    fromWeb("components/**/*.{templ,go,html}"),
    fromWeb("public/**/*.html"),
    fromWeb("styles/**/*.css"),
    fromWeb("../internal/**/*.{templ,go,html}"),
    fromWeb("../domain/**/*.{templ,go,html}"),
    fromWeb("../cmd/**/*.go"),
  ],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dark", "light"],
  },
};
