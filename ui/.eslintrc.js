module.exports = {
  extends: [
    "plugin:perfectionist/recommended-natural-legacy",
    "plugin:@next/next/recommended",
    "plugin:@typescript-eslint/recommended",
  ],
  ignorePatterns: ["components/ui/*.tsx", "generated/**", "electron/**/*.ts"],
  plugins: ["@typescript-eslint/eslint-plugin"],
  rules: {
    "@typescript-eslint/ban-types": "off",
    "@typescript-eslint/no-explicit-any": "off",
  },
};
