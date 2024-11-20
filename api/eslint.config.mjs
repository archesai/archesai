import eslint from "@eslint/js";
import prettierConfig from "eslint-config-prettier";
import jest from "eslint-plugin-jest";
import perfectionist from "eslint-plugin-perfectionist";
import globals from "globals";
import tseslint from "typescript-eslint";

export default [
  {
    ignores: ["dist/", "*/metadata.ts", "node_modules/"],
  },
  {
    files: ["src/**/*.{js,ts,jsx,tsx}", "test/**/*.{js,ts,jsx,tsx}"],
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
    files: ["**/*.spec.js", "**/*.test.js"],
    ...jest.configs["flat/recommended"],
    rules: {
      ...jest.configs["flat/recommended"].rules,
      "jest/prefer-expect-assertions": "off",
    },
  },
  {
    rules: {
      "@typescript-eslint/no-explicit-any": "off",
    },
  },
  perfectionist.configs["recommended-natural"],
  prettierConfig,
];
