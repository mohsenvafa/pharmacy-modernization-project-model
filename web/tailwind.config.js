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
    extend: {
      colors: {
        // map Tailwind color utilities to CSS variables that flip per theme
        "optum-primary": "rgb(var(--optum-primary) / <alpha-value>)",
        "optum-primary-hover":
          "rgb(var(--optum-primary-hover) / <alpha-value>)",
        "optum-outline": "rgb(var(--optum-outline) / <alpha-value>)",
        "optum-muted": "rgb(var(--optum-muted) / <alpha-value>)",
        "optum-cream": "rgb(var(--optum-cream) / <alpha-value>)",
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dark", "light"],
  },
};
