import { z } from 'zod'

export const BaseEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  id: z.ZodString
  updatedAt: z.ZodString
}> = z.object({
  createdAt: z.string().describe('The date this item was created'),
  id: z.string().describe('The ID of the item'),
  updatedAt: z.string().describe('The date this item was last updated')
})

export type BaseEntity = z.infer<typeof BaseEntitySchema>

export type BaseInsertion<TEntity extends BaseEntity> = Omit<TEntity, 'id'> &
  Partial<BaseEntity>
