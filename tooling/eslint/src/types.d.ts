declare module 'eslint-plugin-import' {
  import type { FlatConfig } from '@eslint/compat'
  export const flatConfigs: {
    recommended: FlatConfig
    typescript: FlatConfig
  }
  export const rules: Record<string, Rule.RuleModule>
}
