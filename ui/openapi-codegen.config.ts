import { defineConfig } from '@openapi-codegen/cli'
import {
  generateReactQueryComponents,
  generateSchemaTypes
} from '@openapi-codegen/typescript'

export default defineConfig({
  archesApi: {
    from: {
      relativePath: '../api/openapi-spec.yaml',
      source: 'file'
    },
    outputDir: 'generated',
    to: async (context: any) => {
      const filenamePrefix = 'archesApi'
      const { schemasFiles } = await generateSchemaTypes(context, {
        filenamePrefix
      })
      await generateReactQueryComponents(context, {
        filenamePrefix,
        schemasFiles
      })
    }
  }
})
