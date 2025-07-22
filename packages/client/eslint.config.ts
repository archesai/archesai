import base from '@archesai/eslint/base'
import react from '@archesai/eslint/react'

export default [
  ...react,
  ...base,
  {
    ignores: ['**/generated/**/*.ts']
  }
]
