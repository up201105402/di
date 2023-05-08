import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { fileURLToPath, URL } from "node:url";
import vueJsx from '@vitejs/plugin-vue-jsx';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => ['trix-editor'].indexOf(tag) !== -1
        }
      },
    }),
    vueJsx({
      // options are passed on to @vue/babel-plugin-jsx
    })
  ],
  optimizeDeps: {
    include: [
      'nouislider',
      'wnumb',
      'trix'
    ]
  },
  build: {
    // generate manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: 'src/main.js',
    },
  },
  server: {
    origin: 'http://localhost:8001',
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
})
