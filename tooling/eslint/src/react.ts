import type { ConfigArray } from 'typescript-eslint'

import pluginQuery from '@tanstack/eslint-plugin-query'
import reactPlugin from 'eslint-plugin-react'
import hooksPlugin from 'eslint-plugin-react-hooks'
import globals from 'globals'
import tseslint from 'typescript-eslint'

const react: ConfigArray = tseslint.config({
  files: ['**/*.{ts,tsx}'],
  name: 'react',
  settings: { react: { version: 'detect' } },
  extends: [
    hooksPlugin.configs.recommended,
    reactPlugin.configs.flat['recommended']!,
    reactPlugin.configs.flat['jsx-runtime']!,
    ...pluginQuery.configs['flat/recommended']
  ],
  rules: {
    'react-hooks/react-compiler': 'error',
    'react/prop-types': 'off'
  },
  languageOptions: {
    globals: {
      ...globals.browser,
      ...globals.serviceworker
    }
  }
})

export default react
