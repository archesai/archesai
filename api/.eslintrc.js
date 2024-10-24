module.exports = {
  env: {
    jest: true,
    node: true,
  },
  extends: [
    "plugin:perfectionist/recommended-natural-legacy",
    "plugin:@typescript-eslint/recommended",
    "prettier",
  ],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    project: "tsconfig.json",
    sourceType: "module",
    ecmaVersion: 2020,
  },
  plugins: ["@typescript-eslint", "prettier"],
  root: true,
  ignorePatterns: [".eslintrc.js"],
  rules: {
    "@typescript-eslint/no-explicit-any": "off",
    "@typescript-eslint/interface-name-prefix": "off",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/explicit-module-boundary-types": "off",
  },
};
