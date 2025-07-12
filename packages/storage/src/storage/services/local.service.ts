import { promises } from 'node:fs'
import os from 'node:os'
import { basename, dirname, join } from 'node:path'

import type { HealthCheck, HealthStatus } from '@archesai/core'
import type { FileEntity } from '@archesai/schemas'

import { ConflictException, Logger, NotFoundException } from '@archesai/core'

import { StorageService } from '#storage/storage.service'

/**
 * Service for interacting with the local file system.
 */
export class LocalStorageService extends StorageService implements HealthCheck {
  private readonly health: HealthStatus
  private readonly logger = new Logger(LocalStorageService.name)

  constructor() {
    super()
    this.health = {
      status: 'COMPLETED'
    }
    this.logger.debug('local-storage initialized')
  }

  public async checkExists(
    path: string,
    throwOnMissing = true
  ): Promise<boolean> {
    try {
      await promises.access(path)
      return true
    } catch {
      if (throwOnMissing) {
        throw new NotFoundException(`File at ${path} does not exist`)
      }
      return false
    }
  }

  public async createDirectory(path: string): Promise<void> {
    const exists = await this.checkExists(path, false)
    if (exists) {
      throw new ConflictException(
        'Cannot create directory. File or path already exists at this location'
      )
    }
    await promises.mkdir(path, { recursive: true })
  }

  public async createSignedUrl(
    path: string,
    action: 'read' | 'write'
  ): Promise<FileEntity> {
    if (action === 'write') {
      path = await this.renameIfConflict(path)
      return this.findOne(`file://${path}`)
    } else {
      await this.checkExists(path)
      return this.findOne(`file://${path}`)
    }
  }

  public async delete(path: string): Promise<FileEntity> {
    await this.checkExists(path)
    const stat = await promises.lstat(path)
    if (stat.isDirectory()) {
      await promises.rmdir(path, { recursive: true })
    } else {
      await promises.unlink(path)
    }
    return this.findOne(path)
  }

  public async downloadToBuffer(path: string): Promise<Buffer> {
    await this.checkExists(path)
    const stat = await promises.lstat(path)
    if (stat.isDirectory()) {
      throw new Error('Cannot download a directory')
    }
    return promises.readFile(path)
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
    const stats = await promises.stat(path)
    return {
      createdAt: stats.birthtime.toUTCString(),
      id: path,
      isDir: stats.isDirectory(),
      organizationId: basename(path), // FIXME
      path,
      read: path,
      size: stats.size,
      updatedAt: stats.mtime.toUTCString(),
      write: path
    }
  }

  public getHealth(): HealthStatus {
    return this.health
  }

  public async listDirectory(path: string): Promise<FileEntity[]> {
    await this.checkExists(path)
    const items = await promises.readdir(path)
    return Promise.all(
      items.map(async (item) => {
        const itemFullPath = join(path, item)
        const stats = await promises.stat(itemFullPath)
        const isDir = stats.isDirectory()
        return {
          createdAt: stats.birthtime.toUTCString(),
          id: itemFullPath,
          isDir: isDir,
          name: item,
          organizationId: basename(path), // FIXME
          path: itemFullPath,
          size: stats.size,
          slug: path.split('/').pop() ?? '',
          type: 'file',
          updatedAt: stats.mtime.toUTCString()
        }
      })
    )
  }

  public async uploadFromFile(
    path: string,
    file: {
      buffer: Buffer
    }
  ): Promise<FileEntity> {
    path = await this.renameIfConflict(path)
    await promises.mkdir(dirname(path), { recursive: true })
    await promises.writeFile(path, file.buffer)
    return this.findOne(`file://${path}`)
  }

  public async uploadFromUrl(path: string, url: string): Promise<FileEntity> {
    const response = await fetch(url)
    if (!response.ok) {
      throw new Error(`Failed to fetch file from URL: ${response.statusText}`)
    }
    const buffer = Buffer.from(await response.arrayBuffer())
    const tempFileName = basename(url)
    const tempFilePath = join(os.tmpdir(), tempFileName)
    await promises.writeFile(tempFilePath, buffer)

    try {
      const readUrl = await this.uploadFromFile(path, {
        buffer: buffer
      })
      await promises.unlink(tempFilePath)
      return readUrl
    } catch (err) {
      await promises.unlink(tempFilePath)
      throw err
    }
  }
}
