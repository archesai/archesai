import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinksSchema } from '#http/schemas/links.schema'
import { MetaSchema } from '#http/schemas/meta.schema'
import { ResourceIdentifierSchema } from '#http/schemas/resource-identifier.schema'

export const SuccessDocumentSchema = Type.Object(
  {
    // data: Type.Union([
    //   ResourceObjectSchema,
    //   Type.Array(ResourceObjectSchema),
    //   Type.Null()
    // ]),
    // data: LegacyRef(ResourceObjectSchema),
    included: Type.Optional(Type.Array(LegacyRef(ResourceIdentifierSchema))),
    // jsonapi: Type.Optional(JsonApiObject),
    links: Type.Optional(LegacyRef(LinksSchema)),
    meta: Type.Optional(LegacyRef(MetaSchema))
  },
  {
    $id: 'SuccessDocument',
    description: 'Success Document',
    title: 'Success Document'
  }
)

// included: Type.Optional(
//   Type.Unsafe<StaticDecode<typeof IncludedSchema>>(Type.Ref('Included'))
// ),
// links: Type.Optional(Type.Pick(LinksSchema, ['self'])),
// meta: Type.Optional(
//   Type.Unsafe<StaticDecode<typeof MetaSchema>>(Type.Ref('Meta'))
// )
