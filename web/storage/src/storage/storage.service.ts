import {
  AbortMultipartUploadCommand,
  CompleteMultipartUploadCommand,
  CreateMultipartUploadCommand,
  DeleteObjectCommand,
  GetObjectCommand,
  ListObjectsV2Command,
  PutObjectCommand,
  S3Client,
  UploadPartCommand
} from '@aws-sdk/client-s3'
import { getSignedUrl } from '@aws-sdk/s3-request-presigner'

import type { ConfigService, Logger } from '@archesai/core'

export type StorageService = ReturnType<typeof createStorageService>

interface FileMetadata {
  contentType?: string
  etag: string
  key: string
  lastModified: Date
  size: number
}

interface ListResult {
  continuationToken?: string | undefined
  directories: string[]
  files: FileMetadata[]
}

export function createStorageService(
  configService: ConfigService,
  logger: Logger
) {
  const bucketName = configService.get('storage.bucket')
  const s3 = new S3Client({
    credentials: {
      accessKeyId: configService.get('storage.accesskey'),
      secretAccessKey: configService.get('storage.secretkey')
    },
    endpoint: configService.get('storage.endpoint'),
    forcePathStyle: true,
    region: 'us-east-1'
  })
  logger.debug('s3-storage initialized', { bucketName })

  return {
    async abortMultipartUpload(key: string, uploadId: string): Promise<void> {
      await s3.send(
        new AbortMultipartUploadCommand({
          Bucket: bucketName,
          Key: key,
          UploadId: uploadId
        })
      )
    },

    async completeMultipartUpload(
      key: string,
      uploadId: string,
      parts: { etag: string; partNumber: number }[]
    ): Promise<void> {
      await s3.send(
        new CompleteMultipartUploadCommand({
          Bucket: bucketName,
          Key: key,
          MultipartUpload: {
            Parts: parts.map((part) => ({
              ETag: part.etag,
              PartNumber: part.partNumber
            }))
          },
          UploadId: uploadId
        })
      )
    },

    async createMultipartUpload(
      key: string,
      contentType?: string
    ): Promise<string> {
      const response = await s3.send(
        new CreateMultipartUploadCommand({
          Bucket: bucketName,
          ContentType: contentType,
          Key: key
        })
      )
      if (!response.UploadId) {
        throw new Error('Failed to create multipart upload')
      }
      return response.UploadId
    },

    async delete(key: string): Promise<void> {
      await s3.send(
        new DeleteObjectCommand({
          Bucket: bucketName,
          Key: key
        })
      )
    },

    async download(key: string): Promise<Buffer> {
      const response = await s3.send(
        new GetObjectCommand({
          Bucket: bucketName,
          Key: key
        })
      )

      const chunks: Uint8Array[] = []
      const stream = response.Body as AsyncIterable<Uint8Array>

      for await (const chunk of stream) {
        chunks.push(chunk)
      }

      return Buffer.concat(chunks)
    },

    async getDownloadUrl(key: string, expiresIn = 3600): Promise<string> {
      const command = new GetObjectCommand({
        Bucket: bucketName,
        Key: key
      })
      return getSignedUrl(s3, command, { expiresIn })
    },

    async getMultipartUploadUrl(
      key: string,
      uploadId: string,
      partNumber: number,
      expiresIn = 3600
    ): Promise<string> {
      const command = new UploadPartCommand({
        Bucket: bucketName,
        Key: key,
        PartNumber: partNumber,
        UploadId: uploadId
      })
      return getSignedUrl(s3, command, { expiresIn })
    },

    async getUploadUrl(
      key: string,
      contentType?: string,
      expiresIn = 3600
    ): Promise<string> {
      const command = new PutObjectCommand({
        Bucket: bucketName,
        ContentType: contentType,
        Key: key
      })
      return getSignedUrl(s3, command, { expiresIn })
    },

    async list(prefix?: string, maxKeys = 100): Promise<ListResult> {
      const command = new ListObjectsV2Command({
        Bucket: bucketName,
        Delimiter: '/',
        MaxKeys: maxKeys,
        Prefix: prefix
      })

      const response = await s3.send(command)

      const files: FileMetadata[] = (response.Contents ?? []).map((obj) => ({
        etag: obj.ETag,
        key: obj.Key,
        lastModified: obj.LastModified,
        size: obj.Size
      })) as FileMetadata[]

      const directories = (response.CommonPrefixes ?? [])
        .map((prefix) => prefix.Prefix)
        .filter((f) => f !== undefined)

      return {
        continuationToken: response.NextContinuationToken,
        directories,
        files
      }
    },

    async upload(
      key: string,
      data: Buffer,
      contentType?: string
    ): Promise<void> {
      await s3.send(
        new PutObjectCommand({
          Body: data,
          Bucket: bucketName,
          ContentType: contentType,
          Key: key
        })
      )
    }
  }
}
