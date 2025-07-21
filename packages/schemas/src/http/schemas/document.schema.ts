import { z } from 'zod'

export const DocumentSchemaFactory = <T extends z.ZodType>(
  documentSchema: T
): z.ZodObject<{
  data: T
}> => {
  return z.object({
    data: documentSchema
  })
}

export const DocumentCollectionSchemaFactory = <T extends z.ZodType>(
  documentSchema: T
): z.ZodObject<{
  data: z.ZodArray<T>
  meta: z.ZodObject<{
    total: z.ZodNumber
  }>
}> => {
  return z.object({
    data: z.array(documentSchema),
    meta: z.object({
      total: z.number().describe('Total number of items in the collection')
    })
  })
}

export type DocumentCollectionSchema<T extends z.ZodTypeAny> = ReturnType<
  typeof DocumentCollectionSchemaFactory<T>
>
export type DocumentSchema<T extends z.ZodTypeAny> = ReturnType<
  typeof DocumentSchemaFactory<T>
>
