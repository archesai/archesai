import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { FileEntitySchema } from '@archesai/schemas'

export const CreateFileRequestSchema = Type.Object({
  isDir: FileEntitySchema.properties.isDir,
  path: FileEntitySchema.properties.path
})

export type CreateFileRequest = Static<typeof CreateFileRequestSchema>
