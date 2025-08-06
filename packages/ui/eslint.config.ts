import type { ConfigArray } from '@archesai/eslint'

import { base, react } from '@archesai/eslint'

export default [
  ...react,
  ...base,
  {
    ignores: ['**/shadcn/*']
  }
] as ConfigArray
