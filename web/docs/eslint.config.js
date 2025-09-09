import { base } from "@archesai/eslint/base";
import { react } from "@archesai/eslint/react";

const config = [
  ...react,
  ...base,
  {
    ignores: ["node_modules", "public", "apis", "vite.config.ts"],
  },
];

export default config;
