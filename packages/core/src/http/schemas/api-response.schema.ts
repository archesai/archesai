import type { StaticDecode } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

// import { ErrorsSchema } from '#http/schemas/errors.schema'
import type { IncludedSchema } from '#http/schemas/included.schema'
import type { MetaSchema } from '#http/schemas/meta.schema'

import { LinksSchema } from '#http/schemas/links.schema'

export const ApiResponseSchema = Type.Object(
  {
    included: Type.Optional(
      Type.Unsafe<StaticDecode<typeof IncludedSchema>>(Type.Ref('Included'))
    ),
    links: Type.Optional(Type.Pick(LinksSchema, ['self'])),
    meta: Type.Optional(
      Type.Unsafe<StaticDecode<typeof MetaSchema>>(Type.Ref('Meta'))
    )
  },
  {
    $id: 'ApiResponse',
    description: 'The response object',
    title: 'ApiResponse'
  }
)
