/** @type {import('tailwindcss').Config} */

const withMT = require('@material-tailwind/html/utils/withMT');
module.exports = withMT({
  content: [
    '../view/**/*.html',
    '../view/**/*.templ',
    '../view/**/*.go',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
});
