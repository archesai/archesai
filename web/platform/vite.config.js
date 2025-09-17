import { resolve } from "node:path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import svgr from "vite-plugin-svgr";

export default defineConfig({
  plugins: [
    tanstackStart({
      customViteReactPlugin: true,
      tsr: {
        autoCodeSplitting: true,
        enableRouteTreeFormatting: true,
        routesDirectory: "src/app",
        target: "react",
      },
    }),
    viteReact({
      babel: {
        plugins: [
          [
            "babel-plugin-react-compiler",
            {
              target: "19",
            },
          ],
        ],
      },
    }),
    tailwindcss(),
    svgr(),
  ],
  publicDir: resolve(import.meta.dirname, "../../assets"),
  resolve: {
    alias: {
      "@assets": resolve(import.meta.dirname, "../../assets"),
    },
  },
  server: {
    allowedHosts: ["moose"],
    host: "0.0.0.0",
    port: 3000,
  },
});
