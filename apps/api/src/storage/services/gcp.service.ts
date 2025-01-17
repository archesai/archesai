import { Storage } from '@google-cloud/storage'
import {
  ConflictException,
  Injectable,
  NotFoundException
} from '@nestjs/common'
import axios from 'axios'
import * as path from 'path'

import { StorageItemDto } from '../dto/storage-item.dto'
import { IStorageService } from '../interfaces/storage-provider.interface'
import { v4 } from 'uuid'
import { RunStatusEnum } from '@/src/runs/entities/run.entity'
import { HealthDto } from '@/src/health/dto/health.dto'

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

@Injectable()
export class GoogleCloudStorageService implements IStorageService {
  private readonly bucketName: string
  private readonly expirationTime = 60 * 60 * 1000 // 1 hour in milliseconds
  private readonly storage: Storage

  constructor() {
    this.storage = new Storage({
      credentials: {
        ...archesaiSa
      },
      projectId: 'archesai'
    })
    this.bucketName = 'archesai'
  }

  async checkFileExists(orgname: string, filePath: string): Promise<boolean> {
    const [exists] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .exists()
    return exists
  }

  async createDirectory(orgname: string, dirPath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, dirPath)
    if (exists) {
      throw new ConflictException(
        'Cannot create directory. File or path already exists at this location'
      )
    }
    await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, dirPath) + '/')
      .save('')
  }

  public getHealth(): HealthDto {
    return {
      status: RunStatusEnum.ERROR,
      error: 'Not implemented'
    }
  }

  async delete(orgname: string, filePath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, filePath)
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`)
    }
    await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .delete()
  }

  async download(
    orgname: string,
    filePath: string,
    destination?: string
  ): Promise<{ buffer: Buffer }> {
    const exists = await this.checkFileExists(orgname, filePath)
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`)
    }
    const [buffer] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .download({ destination })
    return { buffer }
  }

  async getMetaData(orgname: string, filePath: string) {
    const exists = await this.checkFileExists(orgname, filePath)
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`)
    }
    const [metadata] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .getMetadata()
    return { metadata }
  }

  async getSignedUrl(
    orgname: string,
    filePath: string,
    action: 'read' | 'write'
  ): Promise<string> {
    let fullPath = this.getFilePath(orgname, filePath)
    if (action === 'write') {
      let conflict = true
      let i = 0
      for (; i < 1000; i++) {
        conflict = await this.checkFileExists(orgname, filePath)
        if (!conflict) {
          break
        }
        filePath = filePath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`)
        fullPath = this.getFilePath(orgname, filePath)
      }
      if (conflict) {
        throw new ConflictException('File already exists')
      }
    } else {
      const exists = await this.checkFileExists(orgname, filePath)
      if (!exists) {
        throw new NotFoundException(`File at ${filePath} does not exist`)
      }
    }

    const [url] = await this.storage
      .bucket(this.bucketName)
      .file(fullPath)
      .getSignedUrl({
        action: action,
        expires: Date.now() + this.expirationTime,
        version: 'v4'
      })

    return url
  }

  async listDirectory(
    orgname: string,
    dirPath: string
  ): Promise<StorageItemDto[]> {
    const fullPath =
      this.getFilePath(orgname, dirPath).replace(/\/+$/, '') + '/'

    const [files] = await this.storage.bucket(this.bucketName).getFiles({
      delimiter: '/',
      prefix: fullPath
    })

    const directories = new Set<string>()
    const fileItems: StorageItemDto[] = []

    files.forEach((file) => {
      const relativePath = file.name.slice(fullPath.length)
      if (relativePath.endsWith('/')) {
        const dirName = relativePath.split('/')[0]
        directories.add(dirName)
      } else if (relativePath) {
        fileItems.push(
          new StorageItemDto({
            createdAt: new Date(file.metadata.timeCreated || Date.now()),
            updatedAt: new Date(file.metadata.updated || Date.now()),
            id: file.id || v4(),
            isDir: false,
            name: relativePath,
            size: Number(file.metadata.size)
          })
        )
      }
    })

    const directoryItems = Array.from(directories).map(
      (dirName) =>
        new StorageItemDto({
          createdAt: new Date(),
          updatedAt: new Date(),
          id: `${fullPath}${dirName}/`,
          isDir: true,
          name: dirName + '/',
          size: 0
        })
    )

    return [...directoryItems, ...fileItems]
  }

  async upload(
    orgname: string,
    filePath: string,
    file: Express.Multer.File
  ): Promise<string> {
    let conflict = await this.checkFileExists(orgname, filePath)
    const originalPath = filePath
    let i = 1
    while (conflict && i < 1000) {
      filePath = originalPath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`)
      conflict = await this.checkFileExists(orgname, filePath)
      i++
    }
    if (conflict) {
      throw new ConflictException('File already exists')
    }

    const ref = this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))

    await ref.save(file.buffer, {
      contentType: file.mimetype,
      metadata: {
        metadata: {
          originalName: file.originalname
        }
      }
    })

    return this.getSignedUrl(orgname, filePath, 'read')
  }

  async uploadFromUrl(
    orgname: string,
    filePath: string,
    url: string
  ): Promise<string> {
    let conflict = await this.checkFileExists(orgname, filePath)
    const originalPath = filePath
    let i = 1
    while (conflict && i < 1000) {
      filePath = originalPath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`)
      conflict = await this.checkFileExists(orgname, filePath)
      i++
    }
    if (conflict) {
      throw new ConflictException('File already exists')
    }

    const ref = this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))

    const response = await axios({
      method: 'get',
      responseType: 'stream',
      url: url
    })

    const writeStream = ref.createWriteStream()

    response.data.pipe(writeStream)

    await new Promise<void>((resolve, reject) => {
      writeStream.on('finish', resolve)
      writeStream.on('error', reject)
    })

    return this.getSignedUrl(orgname, filePath, 'read')
  }

  private getFilePath(orgname: string, filePath: string): string {
    return path.posix.join('storage', orgname, filePath)
  }
}
