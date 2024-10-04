import { defineConfig } from "@openapi-codegen/cli";
import {
  generateReactQueryComponents,
  generateSchemaTypes,
} from "@openapi-codegen/typescript";

export default defineConfig({
  archesApi: {
    from: {
      source: "url",
      url: "http://localhost:3001/-json",
    },
    outputDir: "generated",
    to: async (context: any) => {
      const filenamePrefix = "archesApi";
      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix,
      });
      await generateReactQueryComponents(context, {
        filenamePrefix,
        schemasFiles,
      });
    },
  },
});
