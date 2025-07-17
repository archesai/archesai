import type { Static, TObject, TUnsafe } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { ErrorObjectSchema } from '#http/schemas/error-object.schema'

// import { MetaSchema } from '#http/schemas/meta.schema'

export const ErrorDocumentSchema: TObject<{
  error: TUnsafe<{
    detail: string
    status: string
    title: string
  }>
}> = Type.Object(
  {
    error: LegacyRef(ErrorObjectSchema)
    // meta: Type.Optional(LegacyRef(MetaSchema))
  },
  {
    $id: 'ErrorDocument',
    description: 'Error Document',
    title: 'Error Document'
  }
)

export type ErrorDocument = Static<typeof ErrorDocumentSchema>
