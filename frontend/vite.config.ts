import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
  base: '/',
  plugins: [react()],
  css: {
    modules: {
      localsConvention: "camelCaseOnly",
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8082/",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});