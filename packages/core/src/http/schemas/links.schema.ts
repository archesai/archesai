import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

export const LinkSchema = Type.String({
  $id: 'Link',
  description: 'A link to a related resource',
  format: 'uri',
  title: 'Link'
})

export const LinksSchema = Type.Record(Type.String(), LegacyRef(LinkSchema), {
  $id: 'Links',
  description: 'Collection of links',
  title: 'Links'
})
