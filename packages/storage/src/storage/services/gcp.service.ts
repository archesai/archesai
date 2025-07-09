import { randomUUID } from 'node:crypto'
import { basename } from 'node:path'
import { Readable } from 'node:stream'

import { Storage } from '@google-cloud/storage'

import type { HealthCheck, HealthStatus } from '@archesai/core'
import type { FileEntity } from '@archesai/schemas'

import { ConflictException, Logger, NotFoundException } from '@archesai/core'

import { StorageService } from '#storage/storage.service'

const archesaiSa = {
  auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
  auth_uri: 'https://accounts.google.com/o/oauth2/auth',
  client_email: 'storage@archesai.iam.gserviceaccount.com',
  client_id: '115159155562005966308',
  client_x509_cert_url:
    'https://www.googleapis.com/robot/v1/metadata/x509/storage%40archesai.iam.gserviceaccount.com',
  private_key:
    '-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC8LmWYWLFc6H0P\nND1IT2QhBxhqNc3AsJVt/FPX8DsM1uMe8fkqhC9OSDZ2M50WZnyL3/KIj+C4/Tm2\nRvuA0WCZMdrfcYE+ubQI4qeXCU7cJ+kpgM01iwSUf2yOXp8eQ35onR1hkzPdd/cM\nHtaMeNqgcxIWR4H58aAzk+xUKEn4UFDmtGSCbPEPeKQzaeEDuRQk8oxwF5ggqxC1\nw2lGqK8TjDd9YqVw5y/Gfd4ti4xfP5cizO+TyXUh447Jx4L4tVkhlvM4HfFAskjQ\nwwQfZEHlo3/EuJSwNGaY8rrD/K3o6MXkhEyL6M37AxyUMRuxek1EG5AeWtZmnFc5\nvc54iXOdAgMBAAECggEAAeMQjLxdnES5NaT8n6mAI7P6WDJNlmA5l16U3qZKM/ZB\nqnFFt4SXeYL63NORVdoAPLd/AJnkvdDNj9jas8TQYajg0iXFZIIgiTLwovCVHwWy\nw6zxdyx7hsZXqdPOvP3ImjbycPnfZsSfBbrsxCrVZ7ok/8kxcfaszbzu0pHhhVZl\nvi/k3zKRmATeTfKw7FsCUd6gxWMHBP6Bzvr7khFZCXoFveLgLD0SU1Iafia2XvWW\nOrcVNsMqvXhDgU5IB80rw6sNAgYhoCxZvle4SDfpUsibV4tNQQ3gIsJzCweIH/y7\njvHViuNFEEa9A7DVV2grXBPrRb5XjCrUiKxWFcq+YQKBgQD1I1h23+rHgzXGt4Ij\nlGAUBiQUhPhwVzjQq6p8EAtB0akE1/4xKgi2EiVlB6tN7L1ADixrGzuqEJweFEHv\nFKPK8IKmYaBvrYFm+Uu3YJbpU2EKAQbsby0CiBgs2g1P8f4beD642TQG33HGFTZu\ncrwhB/BXl8bL9ni8qjzpzXaYjQKBgQDEhPpCAHpz2mGDgRnY8fBF1RJj6mvNui2z\nezakejveS6sEdHU/MD6ragiwGvQaN5LSlJpP8Rn3OnCtc7S2VXq1Z35iXSfgYH5e\nz/LvBHTWMMTpJIGgUwcOoveWfgAOv5dItNXH4FtUbo3p1sCe/2/E6F8nhFOFiFbA\nWH9dP9GrUQKBgQDxj+wL8GmGQ2kJsine78ah1M9XHRVIdtr43kE40gKV0IoSyNmn\nDvnYmRcacK1BM8nmRlFFFmf8FTQSe/nhI+CoCctlM40Kn9qFY6JWSStNL6nPVuXA\ntWmQNhZElHdL0XaLEToVo4wePa/690pVGmEC17TiTCFNOksN91/hMWPtvQKBgGek\nBeO3It1kp5bWCE6s0d3SUF+XawFVlfKZIak+ucIzv96amJcZl4OJaUmO/XuyIWGj\nc3qDmgETtgcUBZM/o3Z2PWYc4QHpgdv46ZL6k6++iqq2URK/lvI2KkMY8mjUzDFR\nBYnjHed6YqeXVYDFECoVrtFFbVL4I2BPi+Qe2zHxAoGBAOrtzj75qzqNGSv7fWp3\nAbBClda4Ysd3DGMEBZ85OqVUy3hu+Js1F4znSO3zWu7JjMzy76KRRWlF5YZb/mcr\nXf9lZAjkQQBm1bBiqm/Bpq1npcVKsPLgywFViN8lvA26Vm66qLOLjBzfjY94xVsW\n/UBIFGL9j87Acb2WLMU4OeKs\n-----END PRIVATE KEY-----\n',
  private_key_id: 'c61d819288c05e068ccc53dfe8d7eb1cfb2036f2',
  project_id: 'archesai',
  token_uri: 'https://oauth2.googleapis.com/token',
  type: 'service_account',
  universe_domain: 'googleapis.com'
}

/**
 * Service for interacting with Google Cloud Storage.
 */
export class GoogleCloudStorageService
  extends StorageService
  implements HealthCheck
{
  private readonly bucketName: string
  private readonly expirationTime: number
  private readonly health: HealthStatus
  private readonly logger = new Logger(GoogleCloudStorageService.name)
  private readonly storage: Storage

  constructor() {
    super()
    this.bucketName = 'archesai'
    this.expirationTime = 60 * 60 * 1000
    this.health = {
      status: 'COMPLETED'
    }
    this.storage = new Storage({
      credentials: {
        ...archesaiSa
      },
      projectId: 'archesai'
    })
    this.logger.debug('google-cloud-storage initialized')
  }

  public async checkExists(
    path: string,
    throwOnMissing = true
  ): Promise<boolean> {
    const [exists] = await this.storage
      .bucket(this.bucketName)
      .file(path)
      .exists()
    if (!exists && throwOnMissing) {
      throw new NotFoundException(`File at ${path} does not exist`)
    }
    return exists
  }

  public async createDirectory(path: string): Promise<void> {
    const exists = await this.checkExists(path, false)
    if (exists) {
      throw new ConflictException(
        'Cannot create directory. File or path already exists at this location'
      )
    }
    await this.storage
      .bucket(this.bucketName)
      .file(path + '/')
      .save('')
  }

  public async createSignedUrl(
    path: string,
    action: 'read' | 'write'
  ): Promise<FileEntity> {
    if (action === 'write') {
      path = await this.renameIfConflict(path)
    } else {
      await this.checkExists(path)
    }

    const [url] = await this.storage
      .bucket(this.bucketName)
      .file(path)
      .getSignedUrl({
        action: action,
        expires: Date.now() + this.expirationTime,
        version: 'v4'
      })

    return this.findOne(url)
  }

  public async delete(path: string): Promise<FileEntity> {
    await this.checkExists(path)
    const file = await this.findOne(path)
    await this.storage.bucket(this.bucketName).file(path).delete()
    return file
  }

  public async downloadToBuffer(path: string): Promise<Buffer> {
    await this.checkExists(path)
    const [buffer] = await this.storage
      .bucket(this.bucketName)
      .file(path)
      .download()
    return buffer
  }

  public async downloadToFile(
    path: string,
    destination: string
  ): Promise<void> {
    await this.checkExists(path)
    await this.storage
      .bucket(this.bucketName)
      .file(path)
      .download({ destination })
  }

  public async findOne(path: string): Promise<FileEntity> {
    await this.checkExists(path)
    const [[file]] = await this.storage.bucket(this.bucketName).getFiles({
      prefix: path
    })
    if (!file) {
      throw new NotFoundException(`File at ${path} does not exist`)
    }

    return {
      createdAt: new Date(
        file.metadata.timeCreated ?? Date.now()
      ).toUTCString(),
      id: file.id ?? randomUUID(),
      isDir: false,
      name: path.split('/').pop() ?? '',
      orgname: basename(path), // FIXME
      path: file.name,
      size: Number(file.metadata.size),
      slug: path.split('/').pop() ?? '',
      type: 'file',
      updatedAt: new Date(file.metadata.updated ?? Date.now()).toUTCString()
    }
  }

  public getHealth(): HealthStatus {
    return this.health
  }

  public async listDirectory(path: string): Promise<FileEntity[]> {
    const fullPath = path.replace(/\/+$/, '') + '/'
    const [files] = await this.storage.bucket(this.bucketName).getFiles({
      delimiter: '/',
      prefix: fullPath
    })
    const directories = new Set<string>()
    const fileItems: FileEntity[] = []
    files.forEach((file) => {
      const relativePath = file.name.slice(fullPath.length)
      if (relativePath.endsWith('/')) {
        const dirName = relativePath.split('/')[0] ?? ''
        directories.add(dirName)
      } else if (relativePath) {
        fileItems.push({
          createdAt: new Date(
            file.metadata.timeCreated ?? Date.now()
          ).toUTCString(),
          id: file.id ?? randomUUID(),
          isDir: false,
          name: relativePath,
          orgname: basename(path), // FIXME
          path: file.name,
          size: Number(file.metadata.size),
          slug: path.split('/').pop() ?? '',
          type: 'file',
          updatedAt: new Date(file.metadata.updated ?? Date.now()).toUTCString()
        })
      }
    })
    const directoryItems = Array.from(directories).map(
      (dirName) =>
        ({
          createdAt: new Date().toUTCString(),
          id: `${fullPath}${dirName}/`,
          isDir: true,
          name: dirName + '/',
          orgname: basename(path), // FIXME
          path: `${fullPath}${dirName}/`,
          size: 0,
          slug: path.split('/').pop() ?? '',
          type: 'file',
          updatedAt: new Date().toUTCString()
        }) satisfies FileEntity
    )
    return [...directoryItems, ...fileItems]
  }

  public async uploadFromFile(
    path: string,
    file: {
      buffer: Buffer
      mimetype: string
    }
  ): Promise<FileEntity> {
    path = await this.renameIfConflict(path)
    const ref = this.storage.bucket(this.bucketName).file(path)
    await ref.save(file.buffer, {
      contentType: file.mimetype
    })
    return this.createSignedUrl(path, 'read')
  }

  public async uploadFromUrl(path: string, url: string): Promise<FileEntity> {
    path = await this.renameIfConflict(path)
    const ref = this.storage.bucket(this.bucketName).file(path)
    const response = await fetch(url)
    if (!response.body) {
      throw new Error('Failed to fetch the file from the URL')
    }
    const writeStream = ref.createWriteStream()
    const readable = Readable.fromWeb(response.body)
    readable.pipe(writeStream)
    await new Promise<void>((resolve, reject) => {
      writeStream.on('finish', resolve)
      writeStream.on('error', reject)
    })
    return this.createSignedUrl(path, 'read')
  }
}
