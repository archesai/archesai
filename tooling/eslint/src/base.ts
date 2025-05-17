// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/strict-type-checked.ts
// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/stylistic-type-checked.ts

import path from 'node:path'
import type { ConfigArray } from 'typescript-eslint'

import { includeIgnoreFile } from '@eslint/compat'
import eslint from '@eslint/js'
import prettier from 'eslint-config-prettier'
import importPlugin from 'eslint-plugin-import'
import perfectionist from 'eslint-plugin-perfectionist'
import globals from 'globals'
import tseslint from 'typescript-eslint'

const base: ConfigArray = tseslint.config(
  {
    name: 'ignore .gitignored',
    ...includeIgnoreFile(path.join(import.meta.dirname, '../../../.gitignore'))
  },
  {
    ignores: ['*.config.ts'],
    name: 'ignore config files'
  },
  {
    languageOptions: {
      ecmaVersion: 'latest',
      globals: {
        ...globals.node
      },
      parser: tseslint.parser,
      parserOptions: {
        projectService: true
      },
      sourceType: 'module'
    },
    linterOptions: { reportUnusedDisableDirectives: true },
    name: 'parser',
    plugins: { '@typescript-eslint': tseslint.plugin }
  },
  {
    files: ['**/*.js', '**/*.jsx'],
    name: 'javascript',
    plugins: {
      '@typescript-eslint': tseslint.plugin
    },
    rules: {
      ...eslint.configs.recommended.rules,
      ...tseslint.configs.strictTypeChecked
        .map((c) => c.rules)
        .reduce((a, b) => ({ ...a, ...b }), {}),
      ...tseslint.configs.stylisticTypeChecked
        .map((c) => c.rules)
        .reduce((a, b) => ({ ...a, ...b }), {})
    }
  },
  {
    files: ['**/*.ts', '**/*.tsx'],
    name: 'typescript',
    plugins: {
      '@typescript-eslint': tseslint.plugin
    },
    rules: {
      ...eslint.configs.recommended.rules,
      ...tseslint.configs.strictTypeChecked
        .map((c) => c.rules)
        .reduce((a, b) => ({ ...a, ...b }), {}),
      ...tseslint.configs.stylisticTypeChecked
        .map((c) => c.rules)
        .reduce((a, b) => ({ ...a, ...b }), {}),
      '@typescript-eslint/consistent-type-imports': [
        'warn',
        { fixStyle: 'separate-type-imports', prefer: 'type-imports' }
      ],
      '@typescript-eslint/no-misused-promises': [
        2,
        { checksVoidReturn: { attributes: false } }
      ],
      '@typescript-eslint/no-unnecessary-condition': [
        'error',
        {
          allowConstantLoopConditions: true
        }
      ],
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
      ]
    }
  },
  {
    files: ['**/*.ts', '**/*.tsx'],
    name: 'custom',
    rules: {
      // '@typescript-eslint/consistent-type-assertions': [
      //   'error',
      //   { assertionStyle: 'never' }
      // ],
      '@typescript-eslint/consistent-type-exports': [
        'error',
        { fixMixedExportsWithInlineTypeSpecifier: true }
      ],
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/explicit-member-accessibility': [
        'error',
        {
          accessibility: 'explicit',
          overrides: {
            constructors: 'no-public'
          }
        }
      ],
      '@typescript-eslint/explicit-module-boundary-types': 'off',
      '@typescript-eslint/no-deprecated': 'off',
      '@typescript-eslint/no-extraneous-class': 'off',
      '@typescript-eslint/no-import-type-side-effects': 'error',
      '@typescript-eslint/no-misused-spread': 'off',
      '@typescript-eslint/no-non-null-assertion': 'off'
    }
  },
  {
    name: 'import plugin',
    plugins: {
      import: importPlugin
    },
    rules: {
      ...importPlugin.flatConfigs.recommended.rules,
      ...importPlugin.flatConfigs.typescript.rules,
      'import/consistent-type-specifier-style': ['error', 'prefer-top-level'],
      'import/default': 'off',
      'import/namespace': 'off',
      'import/no-named-as-default-member': 'off',
      'import/no-relative-packages': 'error',
      'import/no-unresolved': 'off',
      ...(process.env.CI === 'true'
        ? {}
        : {
            'import/no-cycle': 'off',
            'import/no-deprecated': 'off',
            'import/no-named-as-default': 'off',
            'import/no-unused-modules': 'off'
          })
      // import/no-extraneous-dependencies https://github.com/import-js/eslint-plugin-import/blob/main/docs/rules/no-extraneous-dependencies.md
    },
    settings: {
      ...importPlugin.flatConfigs.typescript.settings,
      'import/parsers': {
        '@typescript-eslint/parser': ['.ts', '.tsx']
      },
      'import/resolver': {
        node: true,
        typescript: true
      }
    }
  },
  {
    ...perfectionist.configs['recommended-natural'],
    name: 'perfectionist',
    rules: {
      ...perfectionist.configs['recommended-natural'].rules,
      // 'perfectionist/sort-classes': ['error', { newlinesBetween: 'always' }],
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-named-imports': 'off'
    }
  },
  // {
  //   name: 'jsdoc',
  //   ...jsdoc.configs['flat/recommended-typescript']
  // },
  {
    ...prettier,
    name: 'prettier config'
  }
)

export default base
