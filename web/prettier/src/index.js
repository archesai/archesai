/** @typedef {import("prettier").Config} PrettierConfig */
/** @typedef {import("prettier-plugin-tailwindcss").PluginOptions} TailwindConfig */
/** @typedef {import("@ianvs/prettier-plugin-sort-imports").PluginConfig} SortImportsConfig */

/** @type { PrettierConfig | SortImportsConfig | TailwindConfig } */
const prettierConfig = {
  importOrder: [
    "<TYPES>",
    "",
    "^(react/(.*)$)|^(react$)|^(react-native(.*)$)",
    "<THIRD_PARTY_MODULES>",
    "",
    "<TYPES>^@archesai",
    "",
    "^@archesai/(.*)$",
    "",
    "<TYPES>^[.|..|~|#]",
    "",
    "^#",
    "^[../]",
    "^[./]",
  ],
  plugins: [
    "@ianvs/prettier-plugin-sort-imports",
    "prettier-plugin-tailwindcss",
  ],
  singleAttributePerLine: true,
  tailwindFunctions: ["cn", "cva"],
  tailwindStylesheet: "./src/styles/globals.css",
};

export default prettierConfig;
