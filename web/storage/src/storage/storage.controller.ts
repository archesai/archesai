import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import {
  AbortMultipartSchema,
  CompleteMultipartSchema,
  CreateMultipartSchema,
  GetDownloadUrlSchema,
  GetMultipartPartUrlSchema,
  GetUploadUrlSchema,
  KeyParamsSchema,
  ListFilesSchema,
  ListResultSchema,
  SuccessResponseSchema,
  UploadIdResponseSchema,
  UrlResponseSchema
} from '@archesai/schemas'

import type { StorageService } from '#storage/storage.service'

interface StoragePluginOptions {
  storageService: StorageService
}

export const storagePlugin: FastifyPluginAsyncZod<
  StoragePluginOptions
> = async (app, { storageService }) => {
  app.post(
    '/upload-url',
    {
      schema: {
        body: GetUploadUrlSchema,
        response: {
          200: UrlResponseSchema
        }
      }
    },
    async (request) => {
      const url = await storageService.getUploadUrl(
        request.body.key,
        request.body.contentType,
        request.body.expiresIn
      )
      return { url }
    }
  )

  app.post(
    '/download-url',
    {
      schema: {
        body: GetDownloadUrlSchema,
        response: {
          200: UrlResponseSchema
        }
      }
    },
    async (request) => {
      const url = await storageService.getDownloadUrl(
        request.body.key,
        request.body.expiresIn
      )
      return { url }
    }
  )

  app.get(
    '/files',
    {
      schema: {
        querystring: ListFilesSchema,
        response: {
          200: ListResultSchema
        }
      }
    },
    async (request) => {
      return storageService.list(request.query.prefix, request.query.maxKeys)
    }
  )

  app.delete(
    '/files/:key',
    {
      schema: {
        params: KeyParamsSchema,
        response: {
          200: SuccessResponseSchema
        }
      }
    },
    async (request) => {
      await storageService.delete(request.params.key)
      return { success: true }
    }
  )

  app.post(
    '/multipart/create',
    {
      schema: {
        body: CreateMultipartSchema,
        response: {
          200: UploadIdResponseSchema
        }
      }
    },
    async (request) => {
      const uploadId = await storageService.createMultipartUpload(
        request.body.key,
        request.body.contentType
      )
      return { uploadId }
    }
  )

  app.post(
    '/multipart/part-url',
    {
      schema: {
        body: GetMultipartPartUrlSchema,
        response: {
          200: UrlResponseSchema
        }
      }
    },
    async (request) => {
      const url = await storageService.getMultipartUploadUrl(
        request.body.key,
        request.body.uploadId,
        request.body.partNumber,
        request.body.expiresIn
      )
      return { url }
    }
  )

  app.post(
    '/multipart/complete',
    {
      schema: {
        body: CompleteMultipartSchema,
        response: {
          200: SuccessResponseSchema
        }
      }
    },
    async (request) => {
      await storageService.completeMultipartUpload(
        request.body.key,
        request.body.uploadId,
        request.body.parts
      )
      return { success: true }
    }
  )

  app.post(
    '/multipart/abort',
    {
      schema: {
        body: AbortMultipartSchema,
        response: {
          200: SuccessResponseSchema
        }
      }
    },
    async (request) => {
      await storageService.abortMultipartUpload(
        request.body.key,
        request.body.uploadId
      )
      return { success: true }
    }
  )

  await Promise.resolve()
}
