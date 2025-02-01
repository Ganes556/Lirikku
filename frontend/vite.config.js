import inject from '@rollup/plugin-inject';
import tailwindcss from 'tailwindcss';
import autoprefixer from 'autoprefixer';


/** @type {import('vite').UserConfig} */
export default {
  plugins: [
    inject({
      htmx: 'htmx.org',
    }),
  ],
  server: {
    origin: 'http://localhost:5173',
  },
  css: {
    postcss: {
      plugins: [tailwindcss(), autoprefixer()],
    },
  },
  build: {
    manifest: true,
    rollupOptions: {
      input: './main.js',
      output: {
        format: 'iife',
        dir: '../static',
        entryFileNames: 'main.js',
      },
    },
  },
};
