import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      "/archives/": "http://127.0.0.1:8080/"
    }
  },
  plugins: [
    vue({
      template: {
          compilerOptions: {
            // Set replay-web as external element
            isCustomElement: tag => tag.startsWith('replay-web'),
          },
        },
    })
  ],
})
