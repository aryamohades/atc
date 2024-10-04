const defaultTheme = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["static/html/**/*.tmpl"],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Open Sans"', ...defaultTheme.fontFamily.sans],
      },
    },
  },
  daisyui: {
    darkTheme: false,
    themes: [
      {
        corporate: {
          ...require("daisyui/src/theming/themes")["corporate"],
        },
      },
    ],
  },
  plugins: [require("daisyui")],
};
