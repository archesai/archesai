import type { Static, TBoolean, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { FileEntitySchema } from '#storage/files/entities/file.entity'

export const CreateFileDtoSchema: TObject<{
  isDir: TBoolean
  path: TString
}> = Type.Object({
  isDir: FileEntitySchema.properties.isDir,
  path: FileEntitySchema.properties.path
})

export type CreateFileDto = Static<typeof CreateFileDtoSchema>
