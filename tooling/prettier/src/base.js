/** @typedef {import("prettier").Config} PrettierConfig */
/** @typedef {import("prettier-plugin-tailwindcss").PluginOptions} TailwindConfig */
/** @typedef {import("@ianvs/prettier-plugin-sort-imports").PluginConfig} SortImportsConfig */
/** @typedef {import("prettier-plugin-sort-json").SortJsonOptions} SortJsonOptions */

/** @type { PrettierConfig | SortImportsConfig | TailwindConfig | SortJsonOptions} */
export default {
  arrowParens: 'always',
  importOrder: [
    '<TYPES>',
    '',
    '^(react/(.*)$)|^(react$)|^(react-native(.*)$)',
    '^(next/(.*)$)|^(next$)',
    '^(expo(.*)$)|^(expo$)',
    '<THIRD_PARTY_MODULES>',
    '',
    '<TYPES>^@archesai',
    '',
    '^@archesai/(.*)$',
    '',
    '<TYPES>^[.|..|~|#]',
    '',
    '^#',
    '^~/',
    '^[../]',
    '^[./]'
  ],
  importOrderParserPlugins: ['typescript', 'jsx', 'decorators-legacy'],
  importOrderTypeScriptVersion: '5.8.3',
  jsonRecursiveSort: true,
  jsxSingleQuote: true,
  overrides: [
    {
      files: '*.json.hbs',
      options: {
        parser: 'json'
      }
    },
    {
      files: '*.js.hbs',
      options: {
        parser: 'babel'
      }
    }
  ],
  plugins: [
    '@ianvs/prettier-plugin-sort-imports',
    'prettier-plugin-sort-json',
    'prettier-plugin-tailwindcss',
    'prettier-plugin-packagejson',
    'prettier-plugin-sh'
  ],
  printWidth: 80,
  semi: false,
  singleAttributePerLine: true,
  singleQuote: true,
  tabWidth: 2,
  tailwindStylesheet: './src/styles/globals.css',
  tailwindFunctions: ['cn', 'cva'],
  trailingComma: 'none'
}
