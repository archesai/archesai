import { Type } from '@sinclair/typebox'

import { ContentEntitySchema } from '@archesai/domain'

export const CreateContentRequestSchema = Type.Object({
  name: ContentEntitySchema.properties.name,
  text: ContentEntitySchema.properties.text,
  url: ContentEntitySchema.properties.url
})
