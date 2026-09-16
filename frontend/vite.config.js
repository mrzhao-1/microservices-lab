import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 本地开发时，Vite dev server 把 /api 反代到网关，作用等同 nginx 的 location /api/。
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8888',
        changeOrigin: true,
      },
    },
  },
})
