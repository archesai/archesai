import type { CapacitorConfig } from "@capacitor/cli";

const config: CapacitorConfig = {
  appId: "com.archesai.app",
  appName: "App",
  bundledWebRuntime: false,
  includePlugins: ["@capacitor-community/electron"],
  plugins: {
    SplashScreen: {
      launchShowDuration: 0,
    },
  },
  server: {
    cleartext: true,
    url: "http://bob:3000",
  },
  webDir: "out",
};

export default config;
