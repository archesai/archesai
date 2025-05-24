import base from '@archesai/eslint/base'

export default [
  ...base,
  {
    ignores: ['apps/**', 'packages/**', 'tooling/**', 'e2e/**']
  }
]
