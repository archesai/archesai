import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { FileEntitySchema } from '@archesai/domain'

export const CreateSignedUrlRequestSchema = Type.Object({
  action: Type.Union([Type.Literal('read'), Type.Literal('write')], {
    description: 'The type of signed URL to create'
  }),
  path: FileEntitySchema.properties.path
})

export type CreateSignedUrlRequest = Static<typeof CreateSignedUrlRequestSchema>
