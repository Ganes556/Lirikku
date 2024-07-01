import inject from '@rollup/plugin-inject';

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
  build: {
    manifest: true,
    rollupOptions: {
      input: 'main.js',
      output: {
        format: 'iife',
        dir: '../static',
        entryFileNames: 'main.js',
      },
    },
  },
};
