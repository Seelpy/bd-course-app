import { defineConfig } from "vite";
import tsconfigPaths from "vite-tsconfig-paths";
import react from "@vitejs/plugin-react-swc";

// https://vite.dev/config/
export default defineConfig({
  base: "/",
  plugins: [react(), tsconfigPaths()],
  resolve: {
    alias: {
      "@api": "/src/api",
      "@assets": "/src/assets",
      "@shared": "/src/shared",
    },
  },
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
