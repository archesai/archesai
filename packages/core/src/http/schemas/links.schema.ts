import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { MetaSchema } from '#http/schemas/meta.schema'

export const LinkSchema = Type.Union(
  [
    Type.String({ format: 'uri' }),
    Type.Object({
      describedby: Type.Optional(Type.String({ format: 'uri' })),
      href: Type.String({ format: 'uri' }),
      hreflang: Type.Optional(Type.String()),
      meta: Type.Optional(LegacyRef(MetaSchema)),
      rel: Type.Optional(Type.String()),
      title: Type.Optional(Type.String()),
      type: Type.Optional(Type.String())
    })
  ],
  {
    $id: 'Link',
    description: 'Link object or URI string',
    title: 'Link'
  }
)

export const LinksSchema = Type.Object(
  {
    first: Type.Optional(LegacyRef(LinkSchema)),
    last: Type.Optional(LegacyRef(LinkSchema)),
    next: Type.Optional(LegacyRef(LinkSchema)),
    prev: Type.Optional(LegacyRef(LinkSchema)),
    self: LegacyRef(LinkSchema)
  },
  {
    $id: 'Links',
    description: 'Collection of links',
    title: 'Links'
  }
)
