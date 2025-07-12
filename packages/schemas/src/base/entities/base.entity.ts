import type { Static, TSchema } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const BaseEntitySchema = Type.Object({
  createdAt: Type.String({ description: 'The date this item was created' }),
  id: Type.String({ description: 'The ID of the item' }),
  name: Type.String({ description: 'The name of the item' }),
  slug: Type.String({ description: 'The slug of the item' }),
  type: Type.String({ description: 'The type of the item' }),
  updatedAt: Type.String({
    description: 'The date this item was last updated'
  })
})

export abstract class BaseEntity implements Static<typeof BaseEntitySchema> {
  public readonly createdAt: string
  public readonly id: string
  public readonly name: string
  public readonly slug: string
  public abstract type: string
  public readonly updatedAt: string

  protected constructor(props: BaseInsertion<BaseEntity>) {
    this.createdAt = props.createdAt ?? new Date().toISOString()
    this.id = props.id ?? Math.random().toString(36).substring(2, 15)
    // Ensure name is trimmed and not empty
    this.name = props.name?.trim() === '' ? this.id : (props.name ?? this.id)
    this.updatedAt = props.updatedAt ?? this.createdAt
    this.slug = generateSlug(this.name)
  }
}

export const generateSlug = (name: string): string =>
  name
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-') // Replace non-alphanumeric chars with "-"
    .replace(/^-+|-+$/g, '') // Remove leading/trailing dashes

// Extract hard-typed keys
export type BaseEntityKeys = keyof Static<typeof BaseEntitySchema>

export type BaseInsertion<TEntity extends BaseEntity> = Omit<
  TEntity,
  keyof BaseEntity
> &
  Partial<BaseEntity>

export const LegacyRef = <T extends TSchema>(schema: T) =>
  Type.Unsafe<Static<T>>(Type.Ref(schema.$id!))
