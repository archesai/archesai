import type { ConfigArray } from 'typescript-eslint'

import base from '@archesai/eslint/base'

export default [
  ...base,
  {
    ignores: ['**/generated/**/*.ts']
  }
] satisfies ConfigArray
