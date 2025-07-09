import { basename, dirname, extname, join } from 'node:path'

import type { FileEntity } from '@archesai/schemas'

export const STORAGE_SERVICE = 'STORAGE_SERVICE'

export abstract class StorageService {
  public abstract checkExists(
    path: string,
    throwOnMissing: boolean
  ): Promise<boolean>
  public abstract createDirectory(path: string): Promise<void>
  public abstract createSignedUrl(
    path: string,
    action: 'read' | 'write'
  ): Promise<FileEntity>
  public abstract delete(path: string): Promise<FileEntity>
  public abstract downloadToBuffer(path: string): Promise<Buffer>
  public abstract downloadToFile(
    path: string,
    destination: string
  ): Promise<void>
  public abstract findOne(path: string): Promise<FileEntity>
  public abstract listDirectory(path: string): Promise<FileEntity[]>
  public async renameIfConflict(path: string): Promise<string> {
    const originalFilePath = path
    for (let i = 0; i < 10000; i++) {
      const conflict = await this.checkExists(path, false)
      if (!conflict) {
        return path
      }
      const ext = extname(originalFilePath)
      const baseName = basename(originalFilePath, ext)
      const dirName = dirname(originalFilePath)
      path = join(dirName, `${baseName}(${(i + 1).toString()})${ext}`)
    }
    throw new Error('Could not resolve file name conflict')
  }
  public abstract uploadFromFile(
    path: string,
    file: {
      buffer: Buffer
      mimetype: string
      originalname: string
    }
  ): Promise<FileEntity>
  public abstract uploadFromUrl(path: string, url: string): Promise<FileEntity>
}
