import { writeFileSync } from 'fs'
import path from 'path'

import { z } from 'zod'

import { ArchesConfigSchema } from './src/config/config.schema.ts'

const schema = z.toJSONSchema(ArchesConfigSchema, {
  io: 'input'
})

schema.additionalProperties = true

const schemaString = JSON.stringify(schema, null, 2)

const outputPath = path.join(
  import.meta.dirname,
  '../../helm/arches/values.schema.json'
)

writeFileSync(outputPath, schemaString, 'utf-8')

console.log('Schema written to ' + outputPath)
