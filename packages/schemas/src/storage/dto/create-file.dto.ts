import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { FileEntitySchema } from '#storage/entities/file.entity'

export const CreateFileDtoSchema = Type.Object({
  isDir: FileEntitySchema.properties.isDir,
  path: FileEntitySchema.properties.path
})

export type CreateFileDto = Static<typeof CreateFileDtoSchema>
