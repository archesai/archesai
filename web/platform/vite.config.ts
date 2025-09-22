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

  // Only read arches.yaml in development mode
  if (mode === "development") {
    const configPath = resolve(import.meta.dirname, "../../arches.yaml");
    const yamlContent = readFileSync(configPath, "utf8");
    const parsedConfig = yaml.load(yamlContent);

    // Check if config is valid
    if (!parsedConfig || typeof parsedConfig !== "object") {
      throw new Error("Invalid YAML config");
    }

    // Use parsed config directly
    const config = parsedConfig as {
      api?: { host?: string; port?: number; environment?: string };
    };

    // Extract values for environment variables
    const apiHost = config.api?.host || "localhost";
    const apiPort = config.api?.port || 3001;
    const apiUrl = `http://${apiHost}:${apiPort}`;

    // Define environment variables from config
    envVars = {
      "import.meta.env.VITE_API_HOST": JSON.stringify(apiHost),
      "import.meta.env.VITE_API_PORT": JSON.stringify(String(apiPort)),
      "import.meta.env.VITE_API_URL": JSON.stringify(apiUrl),
      "import.meta.env.VITE_ENVIRONMENT": JSON.stringify(
        config.api?.environment || "development",
      ),
    };

    // Set up proxy for API routes
    proxyConfig = {
      "/api": {
        changeOrigin: true,
        secure: false,
        target: apiUrl,
      },
    };
  }

  return {
    define: envVars,
    plugins: [
      tanstackStart({
        customViteReactPlugin: true,
        tsr: {
          // autoCodeSplitting: true,
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
    server: {
      allowedHosts: [],
      host: "0.0.0.0",
      port: 3000,
      proxy: proxyConfig,
      strictPort: true,
    },
  };
});
