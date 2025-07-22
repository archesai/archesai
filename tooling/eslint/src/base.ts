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
    ignores: ['*.config.ts', '*.config.js'],
    name: 'ignore config files'
  },
  // Base JavaScript config
  {
    files: ['**/*.{js,mjs,cjs}'],
    name: 'javascript-base',
    extends: [eslint.configs.recommended],
    languageOptions: {
      ecmaVersion: 'latest',
      globals: globals.node,
      sourceType: 'module'
    },
    linterOptions: { reportUnusedDisableDirectives: true }
  },
  // TypeScript config
  {
    files: ['**/*.{ts,tsx}'],
    name: 'typescript',
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked
    ],
    languageOptions: {
      ecmaVersion: 'latest',
      globals: globals.node,
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: path.join(import.meta.dirname, '../../../')
      },
      sourceType: 'module'
    },
    linterOptions: { reportUnusedDisableDirectives: true },
    rules: {
      'no-restricted-syntax': [
        'error',
        ...banImportExtension('js'),
        ...banImportExtension('jsx'),
        ...banImportExtension('ts'),
        ...banImportExtension('tsx')
      ],
      '@typescript-eslint/consistent-type-imports': [
        'error',
        { fixStyle: 'separate-type-imports', prefer: 'type-imports' }
      ],
      '@typescript-eslint/consistent-type-exports': [
        'error',
        { fixMixedExportsWithInlineTypeSpecifier: false }
      ],
      '@typescript-eslint/no-misused-promises': [
        2,
        { checksVoidReturn: { attributes: false } }
      ],
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
      ],
      '@typescript-eslint/no-import-type-side-effects': 'error',
      // '@typescript-eslint/consistent-type-assertions': [
      //   'warn',
      //   { assertionStyle: 'never' }
      // ],
      // allow while true loops
      '@typescript-eslint/no-unnecessary-condition': [
        'error',
        {
          allowConstantLoopConditions: true
        }
      ]
    }
  },
  // Import plugin config
  {
    files: ['**/*.{ts,tsx}'],
    name: 'import plugin',
    plugins: {
      import: importPlugin
    },
    rules: {
      ...importPlugin.flatConfigs.recommended.rules,
      ...importPlugin.flatConfigs.typescript.rules,
      // 'import/consistent-type-specifier-style': ['error', 'prefer-top-level'],
      'import/default': 'off',
      'import/namespace': 'off',
      'import/no-named-as-default-member': 'off',
      'import/no-relative-packages': 'error',
      'import/no-unresolved': 'off',
      'import/extensions': 'off',
      ...(process.env['CI'] !== 'true' && {
        'import/no-cycle': 'off',
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
        typescript: true,
        node: true
      }
    }
  },
  {
    ...perfectionist.configs['recommended-natural'],
    name: 'perfectionist',
    rules: {
      ...perfectionist.configs['recommended-natural'].rules,
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-named-imports': 'off',
      'perfectionist/sort-decorators': 'off',
      'perfectionist/sort-enums': 'off'
    }
  },
  {
    ...prettier,
    name: 'prettier config'
  }
)

function banImportExtension(extension: string) {
  const message = `Unexpected use of file extension (.${extension}) in import`
  const literalAttributeMatcher = `Literal[value=/\\.${extension}$/]`
  return [
    {
      // import foo from 'bar.js';
      selector: `ImportDeclaration > ${literalAttributeMatcher}.source`,
      message
    },
    {
      // const foo = import('bar.js');
      selector: `ImportExpression > ${literalAttributeMatcher}.source`,
      message
    },
    {
      // type Foo = typeof import('bar.js');
      selector: `TSImportType > TSLiteralType > ${literalAttributeMatcher}`,
      message
    },
    {
      // const foo = require('foo.js');
      selector: `CallExpression[callee.name = "require"] > ${literalAttributeMatcher}.arguments`,
      message
    }
  ]
}

export default base
