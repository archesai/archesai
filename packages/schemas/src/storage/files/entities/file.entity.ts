import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const FileEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  id: z.ZodUUID
  isDir: z.ZodBoolean
  organizationId: z.ZodString
  path: z.ZodString
  read: z.ZodOptional<z.ZodURL>
  size: z.ZodNumber
  updatedAt: z.ZodString
  write: z.ZodOptional<z.ZodString>
}> = BaseEntitySchema.extend({
  isDir: z.boolean().describe('Whether or not this is a directory'),
  organizationId: z.string().describe('The original name of the file'),
  path: z.string().describe('The path to the item'),
  read: z
    .url()
    .optional()
    .describe(
      'The read-only URL that you can use to download the file from secure storage'
    ),
  size: z.number().describe('The size of the item in bytes'),
  write: z
    .string()
    .optional()
    .describe(
      'The write-only URL that you can use to upload the file to secure storage'
    )
}).meta({
  description: 'Schema for File entity',
  id: 'FileEntity'
})

export type FileEntity = z.infer<typeof FileEntitySchema>

export const FILE_ENTITY_KEY = 'files'
