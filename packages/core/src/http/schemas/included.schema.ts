import { Type } from '@sinclair/typebox'

export const IncludedSchema = Type.Array(Type.Unknown(), {
  $id: 'Included',
  description: 'Included related resources',
  title: 'Included'
})
