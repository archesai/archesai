import type { ConfigArray } from 'typescript-eslint'

import base from '@archesai/eslint/base'

export default [
  ...base,
  {
    ignores: ['**/shadcn/*']
  }
] satisfies ConfigArray
