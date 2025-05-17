import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const MetaSchema = Type.Record(Type.String(), Type.Unknown(), {
  $id: 'Meta',
  description: 'Non-standard meta-information',
  title: 'Meta'
})

export type Meta = Static<typeof MetaSchema>
