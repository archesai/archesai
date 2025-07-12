import type { Static, TSchema } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const BaseEntitySchema = Type.Object({
  createdAt: Type.String({ description: 'The date this item was created' }),
  id: Type.String({ description: 'The ID of the item' }),
  updatedAt: Type.String({
    description: 'The date this item was last updated'
  })
})

export type BaseEntity = Static<typeof BaseEntitySchema>

export type BaseInsertion<TEntity extends BaseEntity> = Omit<
  TEntity,
  keyof BaseEntity
> &
  Partial<BaseEntity>

export const LegacyRef = <T extends TSchema>(schema: T) =>
  Type.Unsafe<Static<T>>(Type.Ref(schema.$id!))
