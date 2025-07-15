import type { ConfigArray } from 'typescript-eslint'

import jestPlugin from 'eslint-plugin-jest'
import tseslint from 'typescript-eslint'

const jest: ConfigArray = tseslint.config({
  ...jestPlugin.configs['flat/recommended'],
  files: ['**/*.spec.ts', '**/*.spec.tsx'],
  name: 'jest',
  rules: {
    ...jestPlugin.configs['flat/recommended'].rules,
    // FIXME - these are just to chill on the errors
    '@typescript-eslint/no-unsafe-argument': 'off',
    '@typescript-eslint/no-unsafe-assignment': 'off',
    '@typescript-eslint/no-unsafe-call': 'off',
    '@typescript-eslint/no-unsafe-member-access': 'off',
    '@typescript-eslint/no-unsafe-return': 'off',
    '@typescript-eslint/restrict-plus-operands': 'off',
    '@typescript-eslint/unbound-method': 'off'
  }
})

export default jest
