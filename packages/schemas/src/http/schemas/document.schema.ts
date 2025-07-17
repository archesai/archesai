import type { TArray, TUnsafe } from '@sinclair/typebox'

import type { TObject, TSchema } from '@archesai/schemas'

import { LegacyRef, Type } from '@archesai/schemas'

export const DocumentSchemaFactory = (
  documentSchema: TSchema
): TObject<{
  data: TUnsafe<unknown>
}> => {
  return Type.Object({
    data: LegacyRef(documentSchema)
    //   errors: Type.Optional(Type.Array(LegacyRef(ErrorObjectSchema))),
    //   meta: Type.Optional(LegacyRef(MetaSchema))
  })
}

export const DocumentColectionSchemaFactory = (
  documentSchema: TSchema
  //   metaSchema?: TSchema
): TObject<{
  data: TArray<TUnsafe<unknown>>
}> => {
  return Type.Object({
    data: Type.Array(LegacyRef(documentSchema)),
    meta: Type.Object({
      total: Type.Number({
        description: 'Total number of items in the collection'
      })
    })
    // meta: metaSchema ? LegacyRef(metaSchema) : Type.Optional(Type.Any())
  })
}

export type DocumentCollectionSchema = ReturnType<
  typeof DocumentColectionSchemaFactory
>
export type DocumentSchema = ReturnType<typeof DocumentSchemaFactory>
