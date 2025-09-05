/** @typedef {import("prettier").Config} PrettierConfig */
/** @typedef {import("prettier-plugin-tailwindcss").PluginOptions} TailwindConfig */
/** @typedef {import("@ianvs/prettier-plugin-sort-imports").PluginConfig} SortImportsConfig */

/** @type { PrettierConfig | SortImportsConfig | TailwindConfig } */
const prettierConfig = {
  arrowParens: 'always',
  experimentalTernaries: true,
  importOrder: [
    '<TYPES>',
    '',
    '^(react/(.*)$)|^(react$)|^(react-native(.*)$)',
    '<THIRD_PARTY_MODULES>',
    '',
    '<TYPES>^@archesai',
    '',
    '^@archesai/(.*)$',
    '',
    '<TYPES>^[.|..|~|#]',
    '',
    '^#',
    '^[../]',
    '^[./]'
  ],
  checkIgnorePragma: true,
  importOrderParserPlugins: ['typescript', 'jsx'],
  importOrderTypeScriptVersion: '5.9.2',
  jsxSingleQuote: true,
  plugins: [
    '@ianvs/prettier-plugin-sort-imports',
    'prettier-plugin-tailwindcss'
  ],
  semi: false,
  singleAttributePerLine: true,
  singleQuote: true,
  tailwindFunctions: ['cn', 'cva'],
  tailwindStylesheet: './src/styles/globals.css',
  trailingComma: 'none'
}

export default prettierConfig
