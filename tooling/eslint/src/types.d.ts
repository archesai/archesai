declare module 'eslint-plugin-import' {
  import type { FlatConfig } from '@eslint/compat'
  export const flatConfigs: {
    recommended: FlatConfig
    typescript: FlatConfig
  }
  export const rules: Record<string, Rule.RuleModule>
}

declare module 'eslint-plugin-tailwindcss' {
  import type { ConfigArray } from 'typescript-eslint'
  export const configs: {
    'flat/recommended': ConfigArray
  }
}
