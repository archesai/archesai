// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/strict-type-checked.ts
// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/stylistic-type-checked.ts

import path from "node:path";

import { includeIgnoreFile } from "@eslint/compat";
import eslint from "@eslint/js";
import prettier from "eslint-config-prettier";
import { createTypeScriptImportResolver } from "eslint-import-resolver-typescript";
import { importX } from "eslint-plugin-import-x";
import perfectionist from "eslint-plugin-perfectionist";
import { defineConfig } from "eslint/config";
import globals from "globals";
import tseslint from "typescript-eslint";

const base = defineConfig(
  {
    name: "ignore .gitignored",
    ...includeIgnoreFile(path.join(import.meta.dirname, "../../../.gitignore")),
  },
  // TypeScript config
  {
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked,
    ],
    files: ["**/*.{ts,tsx}"],
    languageOptions: {
      ecmaVersion: "latest",
      globals: globals.node,
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: path.join(import.meta.dirname, "../../../"),
      },
      sourceType: "module",
    },
    linterOptions: { reportUnusedDisableDirectives: true },
    name: "javascript-typescript",
    rules: {
      "@typescript-eslint/consistent-type-assertions": [
        "off",
        { assertionStyle: "never" },
      ],
      "@typescript-eslint/consistent-type-exports": [
        "error",
        { fixMixedExportsWithInlineTypeSpecifier: false },
      ],
      "@typescript-eslint/consistent-type-imports": [
        "error",
        { fixStyle: "separate-type-imports", prefer: "type-imports" },
      ],
      "@typescript-eslint/explicit-module-boundary-types": "error",
      "@typescript-eslint/no-import-type-side-effects": "error",
      "@typescript-eslint/no-misused-promises": [
        2,
        { checksVoidReturn: { attributes: false } },
      ],
      "@typescript-eslint/no-unnecessary-condition": [
        "error",
        {
          allowConstantLoopConditions: true,
        },
      ],
      "@typescript-eslint/no-unused-vars": [
        "error",
        { argsIgnorePattern: "^_", varsIgnorePattern: "^_" },
      ],
      "no-restricted-syntax": [
        "error",
        ...banImportExtension("js"),
        ...banImportExtension("jsx"),
        ...banImportExtension("ts"),
        ...banImportExtension("tsx"),
      ],
    },
  },
  // Import plugin config
  {
    files: ["**/*.{ts,tsx}"],
    name: "import plugin",
    plugins: {
      "import-x": importX,
    },
    rules: {
      "import-x/consistent-type-specifier-style": ["error", "prefer-top-level"],
      "import-x/export": "error",
      "import-x/no-duplicates": "error",
      "import-x/no-relative-packages": "error",
      ...(process.env.CI === "true" && {
        "import-x/no-cycle": "error",
        "import-x/no-deprecated": "error",
        "import-x/no-named-as-default": "error",
        "import-x/no-unused-modules": "error",
      }),
      // import/no-extraneous-dependencies https://github.com/import-js/eslint-plugin-import/blob/main/docs/rules/no-extraneous-dependencies.md
      // 'import-x/enforce-node-protocol-usage': ['error', 'always'],
    },
    settings: {
      "import-x/extensions": [".ts", ".tsx"],
      "import-x/external-module-folders": [
        "node_modules",
        "node_modules/@types",
      ],
      "import-x/parsers": {
        "@typescript-eslint/parser": [".ts", ".tsx"],
      },
      "import-x/resolver-next": [createTypeScriptImportResolver()],
    },
  },
  {
    ...perfectionist.configs["recommended-natural"],
    name: "perfectionist",
    rules: {
      ...perfectionist.configs["recommended-natural"].rules,
      "perfectionist/sort-decorators": "off",
      "perfectionist/sort-enums": "off",
      "perfectionist/sort-imports": "off",
      "perfectionist/sort-maps": "off",
      "perfectionist/sort-named-imports": "off",
      "perfectionist/sort-sets": "off",
    },
  },
  {
    ...prettier,
    name: "prettier config",
  },
);

function banImportExtension(extension) {
  const message = `Unexpected use of file extension (.${extension}) in import`;
  const literalAttributeMatcher = `Literal[value=/\\.${extension}$/]`;
  return [
    {
      message,
      selector: `ImportDeclaration > ${literalAttributeMatcher}.source`, // import foo from 'bar.js';
    },
    {
      message,
      selector: `ImportExpression > ${literalAttributeMatcher}.source`, // const foo = import('bar.js');
    },
    {
      message,
      selector: `TSImportType > TSLiteralType > ${literalAttributeMatcher}`, // type Foo = typeof import('bar.js');
    },
  ];
}

export { base };
