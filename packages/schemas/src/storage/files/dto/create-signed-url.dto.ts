import { z } from 'zod'

import { FileEntitySchema } from '#storage/files/entities/file.entity'

export const CreateSignedUrlDtoSchema: z.ZodObject<{
  action: z.ZodEnum<{
    read: 'read'
    write: 'write'
  }>
  path: z.ZodString
}> = z.object({
  action: z
    .enum(['read', 'write'])
    .describe('The type of signed URL to create'),
  path: FileEntitySchema.shape.path
})

export type CreateSignedUrlDto = z.infer<typeof CreateSignedUrlDtoSchema>
