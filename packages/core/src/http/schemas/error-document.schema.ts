import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { ErrorObjectSchema } from '#http/schemas/error-object.schema'
import { LinksSchema } from '#http/schemas/links.schema'
import { MetaSchema } from '#http/schemas/meta.schema'

export const ErrorDocumentSchema = Type.Object(
  {
    errors: Type.Array(LegacyRef(ErrorObjectSchema)),
    links: Type.Optional(LegacyRef(LinksSchema)),
    meta: Type.Optional(LegacyRef(MetaSchema))
  },
  {
    $id: 'ErrorDocument',
    description: 'Error Document',
    title: 'Error Document'
  }
)
