import type { ConfigExternal } from 'orval'

import { defineConfig } from 'orval'

export default defineConfig({
  archesai: {
    hooks: {
      afterAllFilesWrite: [
        'prettier --write',
        './fix-query-params.sh'
        // `sed -i "s|'\./orval.schemas'|'#generated/orval.schemas'|g" ./src/generated/orval.ts`,
        // `sed -i "s|'../fetcher'|'#fetcher'|g" ./src/generated/orval.ts`
      ]
    },
    input: {
      // validation: {
      //   ...ruleset,
      //   rules: {
      //     ...ruleset.rules,
      //     'ibm-no-unsupported-keywords': {
      //       ...ruleset.rules['ibm-no-unsupported-keywords'],
      //       severity: 'warn'
      //     },
      //     'ibm-property-casing-convention': {
      //       ...ruleset.rules['ibm-property-casing-convention'],
      //       severity: 'off'
      //     },
      //     'ibm-enum-casing-convention': {
      //       ...ruleset.rules['ibm-enum-casing-convention'],
      //       severity: 'off'
      //     },
      //     'ibm-schema-keywords': {
      //       ...ruleset.rules['ibm-schema-keywords'],
      //       severity: 'off'
      //     },
      //     'ibm-path-segment-casing-convention': {
      //       ...ruleset.rules['ibm-path-segment-casing-convention'],
      //       severity: 'off'
      //     }
      //   }
      // },
      target: 'https://api.archesai.dev/docs/openapi.json'
    },
    output: {
      allParamsOptional: true,
      client: 'react-query',
      httpClient: 'fetch',
      mode: 'tags-split',
      override: {
        fetch: {
          includeHttpResponseReturnType: false
        },
        mutator: {
          name: 'customFetch',
          path: './fetcher.ts'
        },
        query: {
          useInfinite: false,
          useQuery: true,
          useSuspenseInfiniteQuery: false,
          useSuspenseQuery: true
        }
      },
      prettier: true,
      // schemas: 'src/generated/models',
      target: './generated/orval.ts',
      urlEncodeParameters: true,
      workspace: './src'
      // propertySortOrder: 'Alphabetical'
    }
  }
} satisfies ConfigExternal)
