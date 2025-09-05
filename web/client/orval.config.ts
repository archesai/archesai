import { defineConfig } from "orval";

export default defineConfig({
  archesai: {
    hooks: {
      afterAllFilesWrite: ["prettier --write"],
    },
    input: {
      target: "../../api/openapi.bundled.yaml",
    },
    output: {
      allParamsOptional: true,
      client: "react-query",
      httpClient: "fetch",
      mode: "tags-split",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          name: "customFetch",
          path: "./fetcher.ts",
        },
        query: {
          useInfinite: false,
          useQuery: true,
          useSuspenseInfiniteQuery: false,
          useSuspenseQuery: true,
        },
      },
      prettier: true,
      // schemas: 'src/generated/models',
      target: "./generated/orval.ts",
      urlEncodeParameters: true,
      workspace: "./src",
      // propertySortOrder: 'Alphabetical'
    },
  },
});
