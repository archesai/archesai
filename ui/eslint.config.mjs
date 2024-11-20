import eslint from "@eslint/js";
import pluginNext from "@next/eslint-plugin-next";
import prettierConfig from "eslint-config-prettier";
import perfectionist from "eslint-plugin-perfectionist";
import globals from "globals";
import tseslint from "typescript-eslint";

export default [
  {
    ignores: ["dist/", ".next/", "node_modules/"],
  },
  {
    files: ["**/*.{js,ts,jsx,tsx}"],
    plugins: {
      "@next/next": pluginNext,
    },
    rules: {
      ...pluginNext.configs.recommended.rules,
      ...pluginNext.configs["core-web-vitals"].rules,
    },
  },
  {
    files: ["**/*.js"],
    languageOptions: {
      sourceType: "commonjs",
    },
  },
  {
    languageOptions: {
      globals: globals.node,
    },
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  {
    rules: {
      "@typescript-eslint/no-explicit-any": "off",
    },
  },
  perfectionist.configs["recommended-natural"],
  prettierConfig,
];
