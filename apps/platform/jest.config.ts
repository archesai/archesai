import type { Config } from 'jest'

import base from '@archesai/jest/base'

export default {
  ...base,
  displayName: 'platform',
  setupFiles: ['@archesai/jest/setup-react'] // FIXME
} satisfies Config
