import type { TSchema } from '@archesai/schemas'

import { LegacyRef, Type } from '@archesai/schemas'

export const DocumentSchemaFactory = (documentSchema: TSchema) => {
  return Type.Object({
    data: LegacyRef(documentSchema)
    //   errors: Type.Optional(Type.Array(LegacyRef(ErrorObjectSchema))),
    //   meta: Type.Optional(LegacyRef(MetaSchema))
  })
}

export const DocumentColectionSchemaFactory = (
  documentSchema: TSchema
  //   metaSchema?: TSchema
) => {
  return Type.Object({
    data: Type.Array(LegacyRef(documentSchema))
    // meta: metaSchema ? LegacyRef(metaSchema) : Type.Optional(Type.Any())
  })
}

export type DocumentCollectionSchema = ReturnType<
  typeof DocumentColectionSchemaFactory
>
export type DocumentSchema = ReturnType<typeof DocumentSchemaFactory>
