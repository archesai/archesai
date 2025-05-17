import type { ConfigArray } from 'typescript-eslint'

import reactPlugin from 'eslint-plugin-react'
// import reactCompiler from 'eslint-plugin-react-compiler'
import hooksPlugin from 'eslint-plugin-react-hooks'

// ...globals.serviceworker,
// ...globals.browswer
const react: ConfigArray = [
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      globals: {
        React: 'writable'
      },
      parserOptions: {
        ...reactPlugin.configs.flat.recommended?.languageOptions.parserOptions,
        ...reactPlugin.configs.flat['jsx-runtime']?.languageOptions
          .parserOptions
      }
    },
    name: 'react',
    plugins: {
      react: reactPlugin,
      'react-hooks': hooksPlugin
    },
    settings: { react: { version: 'detect' } },
    rules: {
      ...reactPlugin.configs.flat.recommended?.rules,
      ...reactPlugin.configs.flat['jsx-runtime']?.rules,
      ...hooksPlugin.configs.recommended.rules
    }
  }
]

export default react
