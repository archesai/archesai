import { resolve } from "node:path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import svgr from "vite-plugin-svgr";

export default defineConfig(({ mode }) => {
  const apiUrl = process.env.VITE_ARCHES_API_HOST || "http://localhost:3001";
  const platformUrl =
    process.env.VITE_ARCHES_PLATFORM_URL || "http://localhost:3000";

  const envVars =
    mode === "development"
      ? {
          "import.meta.env.VITE_ARCHES_API_HOST": JSON.stringify(apiUrl),
          "import.meta.env.VITE_ARCHES_PLATFORM_URL":
            JSON.stringify(platformUrl),
        }
      : {};

  const proxyConfig =
    mode === "development"
      ? {
          "/api": {
            changeOrigin: true,
            rewrite: (path: string) => path.replace(/^\/api/, ""),
            secure: false,
            target: apiUrl,
          },
        }
      : {};

  const allowedHosts =
    mode === "development" ? [platformUrl.replace(/^https?:\/\//, "")] : [];

  return {
    build: {
      outDir: "dist",
    },
    define: envVars,
    plugins: [
      tanstackStart({
        spa: {
          enabled: true,
          prerender: {
            crawlLinks: true,
            outputPath: "index.html",
          },
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
    server: {
      allowedHosts: allowedHosts,
      host: "0.0.0.0",
      port: 3000,
      proxy: proxyConfig,
      strictPort: true,
    },
  };
});
