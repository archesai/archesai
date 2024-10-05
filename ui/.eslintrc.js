module.exports = {
  extends: [
    "plugin:perfectionist/recommended-natural-legacy",
    "plugin:@next/next/recommended",
    "plugin:@typescript-eslint/recommended",
  ],
  ignorePatterns: ["components/ui/*.tsx", "generated/**", "electron/**/*.ts"],
  parser: "@typescript-eslint/parser", // Use TypeScript's parser
  parserOptions: {
    ecmaVersion: 2020, // Allows parsing of modern ECMAScript features
    project: "./tsconfig.json", // Required for rules that need type information
    sourceType: "module", // Allows for the use of imports
  },
  plugins: ["@typescript-eslint"], // Make sure the TypeScript plugin is loaded
  rules: {
    "@typescript-eslint/ban-types": "off",
    "@typescript-eslint/no-explicit-any": "off",
  },
};
