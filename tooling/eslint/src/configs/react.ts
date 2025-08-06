import type { ConfigArray } from 'typescript-eslint'

import pluginQuery from '@tanstack/eslint-plugin-query'
import reactPlugin from 'eslint-plugin-react'
import hooksPlugin from 'eslint-plugin-react-hooks'
import globals from 'globals'
import tseslint from 'typescript-eslint'

const reactFlatConfig = reactPlugin.configs.flat.recommended
const reactJsxRuntimeConfig = reactPlugin.configs.flat['jsx-runtime']
if (!reactFlatConfig || !reactJsxRuntimeConfig) {
  throw new Error(
    'React flat configs are not available. Please check the eslint-plugin-react version.'
  )
}

const react: ConfigArray = tseslint.config({
  extends: [
    hooksPlugin.configs.recommended,
    reactFlatConfig,
    reactJsxRuntimeConfig,
    ...pluginQuery.configs['flat/recommended']
  ],
  files: ['**/*.{ts,tsx}'],
  languageOptions: {
    globals: {
      ...globals.browser,
      ...globals.serviceworker
    }
  },
  name: 'react',
  rules: {
    'react-hooks/react-compiler': 'error',
    'react/prop-types': 'off'
  },
  settings: { react: { version: 'detect' } }
})

export { react }
