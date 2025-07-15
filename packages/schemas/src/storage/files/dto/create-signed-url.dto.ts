import type {
  Static,
  TLiteral,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { FileEntitySchema } from '#storage/files/entities/file.entity'

export const CreateSignedUrlDtoSchema: TObject<{
  action: TUnion<[TLiteral<'read'>, TLiteral<'write'>]>
  path: TString
}> = Type.Object({
  action: Type.Union([Type.Literal('read'), Type.Literal('write')], {
    description: 'The type of signed URL to create'
  }),
  path: FileEntitySchema.properties.path
})

export type CreateSignedUrlDto = Static<typeof CreateSignedUrlDtoSchema>
