import { Type } from '@sinclair/typebox'

export const MetaSchema = Type.Record(Type.String(), Type.Unknown(), {
  $id: 'Meta',
  description: 'Non-standard meta-information',
  title: 'Meta'
})
