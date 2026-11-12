import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { imagetools } from 'vite-imagetools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(), imagetools()],

  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      // เพิ่ม alias นี้เพื่อแก้ปัญหา chunks จำนวนมาก
      '@tabler/icons-react': '@tabler/icons-react/dist/esm/icons/index.mjs',
    },
   
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: [
      'f01.ngrok.dev', // ไม่ต้องใส่ https://
      'localhost',
      '127.0.0.1'
    ],
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})