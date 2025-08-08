// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/strict-type-checked.ts
// https://github.com/typescript-eslint/typescript-eslint/blob/main/packages/typescript-eslint/src/configs/stylistic-type-checked.ts

import path from 'node:path'
import type { ConfigArray } from 'typescript-eslint'

import { includeIgnoreFile } from '@eslint/compat'
import eslint from '@eslint/js'
import nxPlugin from '@nx/eslint-plugin'
import prettier from 'eslint-config-prettier'
import importPlugin from 'eslint-plugin-import'
import perfectionist from 'eslint-plugin-perfectionist'
import globals from 'globals'
import jsoncParser from 'jsonc-eslint-parser'
import tseslint from 'typescript-eslint'

const base: ConfigArray = tseslint.config(
  {
    name: 'ignore .gitignored',
    ...includeIgnoreFile(
      path.join(import.meta.dirname, '../../../../.gitignore')
    )
  },
  // {
  //   ignores: ['*.config.ts', '*.config.js'],
  //   name: 'ignore config files'
  // },
  // Base JavaScript config
  {
    extends: [eslint.configs.recommended],
    files: ['**/*.{js,mjs,cjs}'],
    languageOptions: {
      ecmaVersion: 'latest',
      globals: globals.node,
      sourceType: 'module'
    },
    linterOptions: { reportUnusedDisableDirectives: true },
    name: 'javascript-base'
  },
  // TypeScript config
  {
    extends: [
      eslint.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked
    ],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 'latest',
      globals: globals.node,
      parser: tseslint.parser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: path.join(import.meta.dirname, '../../../../')
      },
      sourceType: 'module'
    },
    linterOptions: { reportUnusedDisableDirectives: true },
    name: 'typescript',
    rules: {
      '@typescript-eslint/consistent-type-exports': [
        'error',
        { fixMixedExportsWithInlineTypeSpecifier: false }
      ],
      '@typescript-eslint/consistent-type-imports': [
        'error',
        { fixStyle: 'separate-type-imports', prefer: 'type-imports' }
      ],
      '@typescript-eslint/no-import-type-side-effects': 'error',
      // '@typescript-eslint/explicit-module-boundary-types': 'error',
      '@typescript-eslint/no-misused-promises': [
        2,
        { checksVoidReturn: { attributes: false } }
      ],
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
      ],
      '@typescript-eslint/no-unused-vars': [
        'error',
        { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }
      ],
      'no-restricted-syntax': [
        'error',
        ...banImportExtension('js'),
        ...banImportExtension('jsx'),
        ...banImportExtension('ts'),
        ...banImportExtension('tsx')
      ]
    }
  },
  // NX config
  {
    files: ['**/*.{ts,tsx}'],
    plugins: { '@nx': nxPlugin },
    rules: {
      '@nx/enforce-module-boundaries': [
        'error',
        {
          allowCircularSelfDependency: true,
          // banTransitiveDependencies: true,
          depConstraints: [
            {
              onlyDependOnLibsWithTags: ['*'],
              sourceTag: '*'
            }
          ]
          // enforceBuildableLibDependency: true
        }
      ]
    }
  },
  {
    files: ['{package,project,nx}.json'],
    languageOptions: {
      parser: jsoncParser
    },
    plugins: { '@nx': nxPlugin },
    rules: {
      '@nx/dependency-checks': [
        'error',
        {
          ignoredDependencies: ['react-dom']
          // includeTransitiveDependencies: true
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
      'import/extensions': 'off',
      'import/namespace': 'off',
      'import/no-named-as-default-member': 'off',
      'import/no-relative-packages': 'error',
      'import/no-unresolved': 'off',
      ...(process.env.CI !== 'true' && {
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
      'perfectionist/sort-decorators': 'off',
      'perfectionist/sort-enums': 'off',
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-named-imports': 'off'
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
      message,
      // import foo from 'bar.js';
      selector: `ImportDeclaration > ${literalAttributeMatcher}.source`
    },
    {
      message,
      // const foo = import('bar.js');
      selector: `ImportExpression > ${literalAttributeMatcher}.source`
    },
    {
      message,
      // type Foo = typeof import('bar.js');
      selector: `TSImportType > TSLiteralType > ${literalAttributeMatcher}`
    },
    {
      message,
      // const foo = require('foo.js');
      selector: `CallExpression[callee.name = "require"] > ${literalAttributeMatcher}.arguments`
    }
  ]
}

export { base }
