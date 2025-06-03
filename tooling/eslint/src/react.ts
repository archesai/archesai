import type { ConfigArray } from 'typescript-eslint'

import reactPlugin from 'eslint-plugin-react'
// import reactCompiler from 'eslint-plugin-react-compiler'
import hooksPlugin from 'eslint-plugin-react-hooks'

// ...globals.serviceworker,
// ...globals.browswer
const react: ConfigArray = [
  {
    files: ['**/*.ts', '**/*.tsx'],
    name: 'react-version',
    settings: { react: { version: 'detect' } }
  },
  {
    files: ['**/*.ts', '**/*.tsx'],
    ...hooksPlugin.configs['recommended-latest']
  },
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      globals: {
        React: 'writable'
      }
    },
    name: 'react',
    ...reactPlugin.configs.flat.recommended,
    ...reactPlugin.configs.flat['jsx-runtime']
  }
]

export default react
