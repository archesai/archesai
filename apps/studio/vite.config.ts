import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import * as yaml from "js-yaml";
import { defineConfig } from "vite";
import svgr from "vite-plugin-svgr";

export default defineConfig(({ mode }) => {
  let envVars = {};
  let proxyConfig = {};
  let allowedHosts = [] as string[];

  // Only read arches.yaml in development mode
  if (mode === "development") {
    console.log("Loading configuration from arches.yaml");
    const configPath = resolve(import.meta.dirname, "./arches.yaml");
    const yamlContent = readFileSync(configPath, "utf8");
    const config = yaml.load(yamlContent) as {
      api?: { url?: string };
      platform?: { url?: string };
    };

    // Check if config is valid
    if (!config || typeof config !== "object") {
      throw new Error("Invalid YAML config");
    }

    const apiUrl = config.api?.url || "http://localhost:3001";
    const platformUrl = config.platform?.url || "http://localhost:3000";

    // Define environment variables from config
    envVars = {
      "import.meta.env.VITE_ARCHES_API_HOST": JSON.stringify(apiUrl),
      "import.meta.env.VITE_ARCHES_PLATFORM_URL": JSON.stringify(platformUrl),
    };

    // Set up proxy for API routes
    proxyConfig = {
      "/api": {
        changeOrigin: true,
        secure: false,
        target: apiUrl,
      },
    };

    allowedHosts = [platformUrl.replace(/^https?:\/\//, "")];
  }

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
