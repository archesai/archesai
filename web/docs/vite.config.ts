import { resolve } from "node:path";

export default {
  publicDir: resolve(import.meta.dirname, "../../assets"),
  server: {
    allowedHosts: ["platform.archesai.dev", "moose"],
  },
};
