import { defineConfig } from 'vite'
import backendPlugin from 'vite-plugin-backend'
import fullReload from 'vite-plugin-full-reload'

export default defineConfig({
  plugins: [
    backendPlugin({
      input: ['resources/assets/js/app.js'],
    }),
    fullReload(['tmp/main']),
  ],
})
