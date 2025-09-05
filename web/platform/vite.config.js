import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
// import { visualizer } from 'rollup-plugin-visualizer'
import { defineConfig } from "vite";

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

export default defineConfig({
  plugins: [
    tailwindcss(),
    tanstackStart({
      customViteReactPlugin: true,
      tsr: {
        enableRouteTreeFormatting: true,
        routesDirectory: "src/app",
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
  ],
  server: {
    allowedHosts: ["platform.archesai.dev"],
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
