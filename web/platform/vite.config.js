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
    allowedHosts: ["platform.archesai.dev", "moose"],
    host: "0.0.0.0",
    port: 3000,
  },
  test: {
    coverage: {
      provider: "v8",
      reportsDirectory: ".coverage",
    },
    environment: "node",
    globals: true,
    include: ["src/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    reporters: ["default"],
    watch: false,
  },
});

// const clientViz = visualizer({
//   brotliSize: true,
//   emitFile: true,
//   filename: 'stats-client.html',
//   template: 'treemap'
// })

// const ssrViz = visualizer({
//   brotliSize: true,
//   emitFile: true,
//   filename: 'stats-ssr.html',
//   template: 'treemap'
// })

// import { visualizer } from 'rollup-plugin-visualizer'

// client-only visualizer
// {
//   ...clientViz,
//   apply: (_c, env) => env.command === 'build' && !env.isSsrBuild
// },
// ssr-only visualizer
// {
//   ...ssrViz,
//   apply: (_c, env) => env.command === 'build' && env.isSsrBuild === true
// }
