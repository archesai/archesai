import type { ConfigArray } from '@archesai/eslint'

import { base, react } from '@archesai/eslint'

const config: ConfigArray = [
  ...react,
  ...base,
  {
    ignores: ['**/shadcn/*']
  }
]

export default config
