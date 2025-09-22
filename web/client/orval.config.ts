import { defineConfig } from "orval";

export default defineConfig({
  archesai: {
    input: {
      target: "../../api/openapi.bundled.yaml",
    },
    output: {
      allParamsOptional: true,
      biome: true,
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
          version: 5,
        },
      },
      // schemas: 'src/generated/models',
      target: "./generated/orval.ts",
      urlEncodeParameters: true,
      workspace: "./src",
      // propertySortOrder: 'Alphabetical'
    },
  },
  archesaiZod: {
    input: {
      target: "../../api/openapi.bundled.yaml",
    },
    output: {
      biome: true,
      client: "zod",
      mode: "single",
      target: "./generated/zod.ts",
      workspace: "./src",
    },
  },
});
