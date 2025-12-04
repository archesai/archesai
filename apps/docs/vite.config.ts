import { resolve } from "node:path";
import { defineConfig } from "vite";
import svgr from "vite-plugin-svgr";

export default defineConfig({
  optimizeDeps: {
    include: ["@archesai/ui"],
  },
  plugins: [svgr()],
  publicDir: resolve(import.meta.dirname, "../../assets"),
  server: {
    allowedHosts: ["moose"],
    host: "0.0.0.0",
    port: 3002,
  },
});
