import type { z } from 'zod'

import { FileEntitySchema } from '#storage/files/entities/file.entity'

export const CreateFileDtoSchema: z.ZodObject<{
  isDir: z.ZodBoolean
  path: z.ZodString
}> = FileEntitySchema.pick({
  isDir: true,
  path: true
})

export type CreateFileDto = z.infer<typeof CreateFileDtoSchema>
