import type { ConfigArray } from 'typescript-eslint'

import nextPlugin from '@next/eslint-plugin-next'

// ...globals.serviceworker,
export default [
  {
    files: ['**/*.ts', '**/*.tsx'],
    name: 'nextjs',
    plugins: {
      '@next/next': nextPlugin
    },
    rules: {
      ...nextPlugin.configs.recommended.rules,
      ...nextPlugin.configs['core-web-vitals'].rules
    }
  }
] satisfies ConfigArray
