/** @typedef {import("prettier").Config} PrettierConfig */
/** @typedef {import("prettier-plugin-tailwindcss").PluginOptions} TailwindConfig */
/** @typedef {import("@ianvs/prettier-plugin-sort-imports").PluginConfig} SortImportsConfig */
/** @typedef {import("prettier-plugin-sort-json").SortJsonOptions} SortJsonOptions */

/** @type { PrettierConfig | SortImportsConfig | TailwindConfig | SortJsonOptions} */
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
  importOrderParserPlugins: ['typescript', 'jsx'],
  importOrderTypeScriptVersion: '5.9.2',
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
  tailwindFunctions: ['cn', 'cva'],
  tailwindStylesheet: './src/styles/globals.css',
  trailingComma: 'none'
}

export default prettierConfig
