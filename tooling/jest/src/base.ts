import type { Config } from 'jest'

import swcrc from '@archesai/swc/base'

export default {
  coverageDirectory: '<rootDir>/coverage',
  moduleFileExtensions: ['ts', 'js', 'html'],
  rootDir: 'src',
  testEnvironment: 'node',
  // The regexp pattern or array of patterns that Jest uses to detect test files
  testRegex: '.*\\.spec\\.ts$',
  transform: {
    '^.+\\.[tj]s$': ['@swc/jest', swcrc]
  }
} satisfies Config
