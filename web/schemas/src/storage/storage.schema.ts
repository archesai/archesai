import { z } from 'zod'

// Response schemas
export const UrlResponseSchema: z.ZodObject<{
  url: z.ZodString
}> = z.object({ url: z.string() })
export const SuccessResponseSchema: z.ZodObject<{
  success: z.ZodBoolean
}> = z.object({ success: z.boolean() })
export const UploadIdResponseSchema: z.ZodObject<{
  uploadId: z.ZodString
}> = z.object({ uploadId: z.string() })

export const FileMetadataSchema: z.ZodObject<{
  contentType: z.ZodOptional<z.ZodString>
  etag: z.ZodString
  key: z.ZodString
  lastModified: z.ZodDate
  size: z.ZodNumber
}> = z.object({
  contentType: z.string().optional(),
  etag: z.string(),
  key: z.string(),
  lastModified: z.date(),
  size: z.number()
})

export const ListResultSchema: z.ZodObject<{
  continuationToken: z.ZodOptional<z.ZodString>
  directories: z.ZodArray<z.ZodString>
  files: z.ZodArray<
    z.ZodObject<{
      contentType: z.ZodOptional<z.ZodString>
      etag: z.ZodString
      key: z.ZodString
      lastModified: z.ZodDate
      size: z.ZodNumber
    }>
  >
}> = z.object({
  continuationToken: z.string().optional(),
  directories: z.array(z.string()),
  files: z.array(FileMetadataSchema)
})

// Request schemas
export const GetUploadUrlSchema: z.ZodObject<{
  contentType: z.ZodOptional<z.ZodString>
  expiresIn: z.ZodDefault<z.ZodNumber>
  key: z.ZodString
}> = z.object({
  contentType: z.string().optional(),
  expiresIn: z.number().default(3600),
  key: z.string()
})

export const GetDownloadUrlSchema: z.ZodObject<{
  expiresIn: z.ZodDefault<z.ZodNumber>
  key: z.ZodString
}> = z.object({
  expiresIn: z.number().default(3600),
  key: z.string()
})

export const KeyParamsSchema: z.ZodObject<{
  key: z.ZodString
}> = z.object({
  key: z.string()
})

export const ListFilesSchema: z.ZodObject<{
  maxKeys: z.ZodDefault<z.ZodNumber>
  prefix: z.ZodOptional<z.ZodString>
}> = z.object({
  maxKeys: z.number().default(100),
  prefix: z.string().optional()
})

export const CreateMultipartSchema: z.ZodObject<{
  contentType: z.ZodOptional<z.ZodString>
  key: z.ZodString
}> = z.object({
  contentType: z.string().optional(),
  key: z.string()
})

export const GetMultipartPartUrlSchema: z.ZodObject<{
  expiresIn: z.ZodDefault<z.ZodNumber>
  key: z.ZodString
  partNumber: z.ZodNumber
  uploadId: z.ZodString
}> = z.object({
  expiresIn: z.number().default(3600),
  key: z.string(),
  partNumber: z.number(),
  uploadId: z.string()
})

export const CompleteMultipartSchema: z.ZodObject<{
  key: z.ZodString
  parts: z.ZodArray<
    z.ZodObject<{
      etag: z.ZodString
      partNumber: z.ZodNumber
    }>
  >
  uploadId: z.ZodString
}> = z.object({
  key: z.string(),
  parts: z.array(
    z.object({
      etag: z.string(),
      partNumber: z.number()
    })
  ),
  uploadId: z.string()
})

export const AbortMultipartSchema: z.ZodObject<{
  key: z.ZodString
  uploadId: z.ZodString
}> = z.object({
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
