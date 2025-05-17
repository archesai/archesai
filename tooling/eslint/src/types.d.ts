declare module 'eslint-plugin-import' {
  import type { FlatConfig } from '@eslint/compat'
  export const flatConfigs: {
    recommended: FlatConfig
    typescript: FlatConfig
  }
  export const rules: Record<string, Rule.RuleModule>
}

declare module '@next/eslint-plugin-next' {
  import type { Linter, Rule } from 'eslint'

  export const configs: {
    'core-web-vitals': { rules: Linter.RulesRecord }
    recommended: { rules: Linter.RulesRecord }
  }
  export const rules: Record<string, Rule.RuleModule>
}

declare module 'eslint-plugin-tailwindcss' {
  import type { ConfigArray } from 'typescript-eslint'
  export const configs: {
    'flat/recommended': ConfigArray
  }
}
