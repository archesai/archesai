import { z } from 'zod'

// Response schemas
export const UrlResponseSchema = z.object({ url: z.string() })
export const SuccessResponseSchema = z.object({ success: z.boolean() })
export const UploadIdResponseSchema = z.object({ uploadId: z.string() })

export const FileMetadataSchema = z.object({
  contentType: z.string().optional(),
  etag: z.string(),
  key: z.string(),
  lastModified: z.date(),
  size: z.number()
})

export const ListResultSchema = z.object({
  continuationToken: z.string().optional(),
  directories: z.array(z.string()),
  files: z.array(FileMetadataSchema)
})

// Request schemas
export const GetUploadUrlSchema = z.object({
  contentType: z.string().optional(),
  expiresIn: z.number().default(3600),
  key: z.string()
})

export const GetDownloadUrlSchema = z.object({
  expiresIn: z.number().default(3600),
  key: z.string()
})

export const KeyParamsSchema = z.object({
  key: z.string()
})

export const ListFilesSchema = z.object({
  maxKeys: z.number().default(100),
  prefix: z.string().optional()
})

export const CreateMultipartSchema = z.object({
  contentType: z.string().optional(),
  key: z.string()
})

export const GetMultipartPartUrlSchema = z.object({
  expiresIn: z.number().default(3600),
  key: z.string(),
  partNumber: z.number(),
  uploadId: z.string()
})

export const CompleteMultipartSchema = z.object({
  key: z.string(),
  parts: z.array(
    z.object({
      etag: z.string(),
      partNumber: z.number()
    })
  ),
  uploadId: z.string()
})

export const AbortMultipartSchema = z.object({
  key: z.string(),
  uploadId: z.string()
})

export type AbortMultipart = z.infer<typeof AbortMultipartSchema>
export type CompleteMultipart = z.infer<typeof CompleteMultipartSchema>
export type CreateMultipart = z.infer<typeof CreateMultipartSchema>
export type FileMetadata = z.infer<typeof FileMetadataSchema>
export type GetDownloadUrl = z.infer<typeof GetDownloadUrlSchema>
export type GetMultipartPartUrl = z.infer<typeof GetMultipartPartUrlSchema>
export type GetUploadUrl = z.infer<typeof GetUploadUrlSchema>
export type KeyParams = z.infer<typeof KeyParamsSchema>
export type ListFiles = z.infer<typeof ListFilesSchema>
export type ListResult = z.infer<typeof ListResultSchema>
export type SuccessResponse = z.infer<typeof SuccessResponseSchema>
export type UploadIdResponse = z.infer<typeof UploadIdResponseSchema>
// Types
export type UrlResponse = z.infer<typeof UrlResponseSchema>
