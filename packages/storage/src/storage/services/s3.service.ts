import { randomUUID } from 'node:crypto'
import { promises } from 'node:fs'
import { basename } from 'node:path'
import { Readable } from 'node:stream'

import {
  // CreateBucketCommand,
  DeleteObjectCommand,
  GetObjectCommand,
  HeadObjectCommand,
  ListObjectsV2Command,
  PutObjectCommand,
  S3Client
} from '@aws-sdk/client-s3'
import { getSignedUrl } from '@aws-sdk/s3-request-presigner'

import type { ConfigService, Logger } from '@archesai/core'
import type { FileEntity } from '@archesai/schemas'

import {
  ConflictException,
  InternalServerErrorException,
  NotFoundException
} from '@archesai/core'

import { StorageService } from '#storage/storage.service'

/**
 * Service for interacting with the S3 storage service.
 */
export class S3StorageProvider extends StorageService {
  private readonly bucketName: string
  private readonly expirationTime = 60 * 60 * 1000
  private readonly s3Client: S3Client

  constructor(configService: ConfigService, logger: Logger) {
    super()
    this.bucketName = configService.get('storage.bucket')

    this.s3Client = new S3Client({
      credentials: {
        accessKeyId: configService.get('storage.accesskey'),
        secretAccessKey: configService.get('storage.secretkey')
      },
      endpoint: configService.get('storage.endpoint'),
      forcePathStyle: true,
      region: 'us-east-1'
    })
    logger.debug('s3-storage initialized')
  }

  public async checkExists(
    path: string,
    throwOnMissing = true
  ): Promise<boolean> {
    try {
      await this.s3Client.send(
        new HeadObjectCommand({
          Bucket: this.bucketName,
          Key: path
        })
      )
      return true
    } catch (error: unknown) {
      if (error instanceof Error && error.name === 'NotFound') {
        if (throwOnMissing) {
          throw new NotFoundException(`File at ${path} does not exist`)
        }
        return false
      }
      throw error
    }
  }

  public async createDirectory(path: string): Promise<void> {
    const exists = await this.checkExists(path, false)
    if (exists) {
      throw new ConflictException(
        'Cannot create directory. File or path already exists at this location'
      )
    }
    const key = path.replace(/\/?$/, '/') + '/'
    await this.s3Client.send(
      new PutObjectCommand({
        Body: '',
        Bucket: this.bucketName,
        Key: key
      })
    )
  }

  public async createSignedUrl(
    path: string,
    action: 'read' | 'write'
  ): Promise<FileEntity> {
    let command
    if (action === 'read') {
      await this.checkExists(path)
      command = new GetObjectCommand({
        Bucket: this.bucketName,
        Key: path
      })
    } else {
      path = await this.renameIfConflict(path)
      command = new PutObjectCommand({
        Bucket: this.bucketName,
        Key: path
      })
    }
    const signedUrl = await getSignedUrl(this.s3Client, command, {
      expiresIn: this.expirationTime / 1000
    })
    return this.findOne(signedUrl)
  }

  public async delete(path: string): Promise<FileEntity> {
    await this.checkExists(path)
    const file = await this.findOne(path)
    await this.s3Client.send(
      new DeleteObjectCommand({
        Bucket: this.bucketName,
        Key: path
      })
    )
    return file
  }

  public async downloadToBuffer(path: string): Promise<Buffer> {
    await this.checkExists(path)
    const result = await this.s3Client.send(
      new GetObjectCommand({
        Bucket: this.bucketName,
        Key: path
      })
    )
    const stream = result.Body
    if (!stream || !(stream instanceof Readable)) {
      throw new InternalServerErrorException('Failed to read file stream')
    }
    const chunks = []
    for await (const chunk of stream) {
      chunks.push(chunk)
    }
    return Buffer.concat(chunks)
  }

  public async downloadToFile(
    path: string,
    destination: string
  ): Promise<void> {
    const buffer = await this.downloadToBuffer(path)
    await promises.writeFile(destination, buffer)
  }

  public async findOne(path: string): Promise<FileEntity> {
    await this.checkExists(path)
    const result = await this.s3Client.send(
      new HeadObjectCommand({
        Bucket: this.bucketName,
        Key: path
      })
    )
    return {
      createdAt: (result.LastModified ?? new Date()).toUTCString(),
      id: path,
      isDir: false,
      organizationId: basename(path), // FIXME
      path,
      size: result.ContentLength ?? 0,
      updatedAt: (result.LastModified ?? new Date()).toUTCString()
    }
  }

  public async listDirectory(path: string): Promise<FileEntity[]> {
    const prefix = path.replace(/\/?$/, '/') + '/'
    const result = await this.s3Client.send(
      new ListObjectsV2Command({
        Bucket: this.bucketName,
        Delimiter: '/',
        Prefix: prefix
      })
    )

    const items: FileEntity[] = []

    if (result.CommonPrefixes) {
      for (const commonPrefix of result.CommonPrefixes) {
        const subPrefix = commonPrefix.Prefix ?? ''
        items.push({
          createdAt: new Date().toUTCString(),
          id: subPrefix,
          isDir: true,
          organizationId: basename(path), // FIXME
          path: subPrefix,
          size: 0,
          updatedAt: new Date().toUTCString()
        })
      }
    }

    if (result.Contents) {
      for (const content of result.Contents) {
        if (content.Key === prefix) {
          continue // Skip the directory placeholder object
        }
        items.push({
          createdAt: (content.LastModified ?? new Date()).toUTCString(),
          id: content.Key ?? randomUUID(),
          isDir: false,
          organizationId: basename(path), // FIXME
          path: content.Key ?? '',
          size: content.Size ?? 0,
          updatedAt: (content.LastModified ?? new Date()).toUTCString()
        })
      }
    }

    return items
  }

  public async uploadFromFile(
    path: string,
    file: {
      buffer: Buffer
      mimetype: string
    }
  ): Promise<FileEntity> {
    path = await this.renameIfConflict(path)
    await this.s3Client.send(
      new PutObjectCommand({
        Body: file.buffer,
        Bucket: this.bucketName,
        ContentType: file.mimetype,
        Key: path
      })
    )
    return this.createSignedUrl(path, 'read')
  }

  public async uploadFromUrl(path: string, url: string): Promise<FileEntity> {
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`Failed to fetch file from URL: ${response.statusText}`)
    }

    const contentType =
      response.headers.get('Content-Type') ?? 'application/octet-stream'
    const buffer = Buffer.from(await response.arrayBuffer())

    return this.uploadFromFile(path, {
      buffer,
      mimetype: contentType
    })
  }

  // private async createBucketIfNotExists() {
  //   try {
  //     await this.s3Client.send(
  //       new CreateBucketCommand({ Bucket: this.bucketName })
  //     )
  //     this.logger.debug(
  //       `Bucket '${this.bucketName}' created or already exists.`
  //     )
  //   } catch (error: unknown) {
  //     if (error instanceof Error && error.name !== 'BucketAlreadyOwnedByYou') {
  //       throw error
  //     }
  //     this.logger.debug(`Bucket '${this.bucketName}' already owned by you.`)
  //   }
  // }
}
