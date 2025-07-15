import type { TArray, TObject, TUnsafe } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { ErrorObjectSchema } from '#http/schemas/error-object.schema'

// import { MetaSchema } from '#http/schemas/meta.schema'

export const ErrorDocumentSchema: TObject<{
  errors: TArray<
    TUnsafe<{
      detail: string
      status: string
      title: string
    }>
  >
}> = Type.Object(
  {
    errors: Type.Array(LegacyRef(ErrorObjectSchema))
    // meta: Type.Optional(LegacyRef(MetaSchema))
  },
  {
    $id: 'ErrorDocument',
    description: 'Error Document',
    title: 'Error Document'
  }
)
