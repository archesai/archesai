import { ConfigExternal, defineConfig } from 'orval'

export default defineConfig({
  archesai: {
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
      override: {
        fetch: {
          includeHttpResponseReturnType: false
        },
        mutator: {
          path: './fetcher.ts',
          name: 'customFetch'
        },
        query: {
          useQuery: true,
          useInfinite: false,
          useSuspenseQuery: true,
          useSuspenseInfiniteQuery: false
        }
      },
      // schemas: 'src/generated/models',
      target: './generated/orval.ts',
      client: 'react-query',
      httpClient: 'fetch',
      mode: 'tags-split',
      urlEncodeParameters: true,
      allParamsOptional: true,
      workspace: './src'
      // propertySortOrder: 'Alphabetical'
    },
    hooks: {
      afterAllFilesWrite: [
        'prettier --write'
        // `sed -i "s|'\./orval.schemas'|'#generated/orval.schemas'|g" ./src/generated/orval.ts`,
        // `sed -i "s|'../fetcher'|'#fetcher'|g" ./src/generated/orval.ts`
      ]
    }
  }
} satisfies ConfigExternal)
