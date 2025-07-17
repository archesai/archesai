import type {
  Static,
  TObject,
  TSchema,
  TString,
  TUnsafe
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const BaseEntitySchema: TObject<{
  createdAt: TString
  id: TString
  updatedAt: TString
}> = Type.Object({
  createdAt: Type.String({ description: 'The date this item was created' }),
  id: Type.String({ description: 'The ID of the item' }),
  updatedAt: Type.String({
    description: 'The date this item was last updated'
  })
})

export type BaseEntity = Static<typeof BaseEntitySchema>

export type BaseInsertion<TEntity extends BaseEntity> = Omit<TEntity, 'id'> &
  Partial<BaseEntity>

export const LegacyRef = <T extends TSchema>(
  schema: T
): TUnsafe<
  (T & {
    params: []
  })['static']
> => {
  if (!schema.$id) {
    throw new Error('Schema must have an $id property')
  }
  return Type.Unsafe<Static<T>>(Type.Ref(schema.$id))
}
