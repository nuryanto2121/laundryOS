import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import { fileURLToPath, URL } from 'node:url';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/auth': 'http://localhost:8080',
      '/master': 'http://localhost:8080',
      '/customers': 'http://localhost:8080',
      '/orders': 'http://localhost:8080',
      '/expenses': 'http://localhost:8080',
      '/reports': 'http://localhost:8080',
      '/attendances': 'http://localhost:8080',
      '/leaves': 'http://localhost:8080',
    },
  },
});
