import type {
  BaseService,
  SearchQuery,
  WebsocketsService
} from '@archesai/core'
import type { BaseInsertion, FileEntity } from '@archesai/schemas'

import { FILE_ENTITY_KEY } from '@archesai/schemas'

import type { StorageService } from '#storage/storage.service'

/**
 * Service for files.
 */
export class FilesService
  implements Omit<BaseService<FileEntity>, 'repository'>
{
  private readonly storageService: StorageService
  private readonly websocketsService: WebsocketsService

  constructor(
    storageService: StorageService,
    websocketsService: WebsocketsService
  ) {
    this.storageService = storageService
    this.websocketsService = websocketsService
  }

  public async create(value: BaseInsertion<FileEntity>): Promise<FileEntity> {
    return this.storageService.uploadFromUrl(value.path, value.path) // FIXME
  }

  public async createMany(
    values: BaseInsertion<FileEntity>[]
  ): Promise<{ count: number; data: FileEntity[] }> {
    const entities = await Promise.all(
      values.map(
        (value) => this.storageService.uploadFromUrl(value.path, value.path) // FIXME
      )
    )
    return { count: entities.length, data: entities }
  }

  public delete(id: string): Promise<FileEntity> {
    return this.storageService.delete(id)
  }

  public async deleteMany(
    query: SearchQuery<FileEntity>
  ): Promise<{ count: number; data: FileEntity[] }> {
    const toDelete = await this.findMany(query)
    const deleted = await Promise.all(
      toDelete.data.map((entity) => this.delete(entity.id))
    )
    return { count: deleted.length, data: deleted }
  }

  public findMany(
    _query: SearchQuery<FileEntity>
  ): Promise<{ count: number; data: FileEntity[] }> {
    throw new Error('Method not implemented.')
  }

  public findOne(id: string): Promise<FileEntity> {
    return this.storageService.findOne(id)
  }

  public update(
    _id: string,
    _data: Partial<BaseInsertion<FileEntity>>
  ): Promise<FileEntity> {
    throw new Error('Method not implemented.')
  }

  public updateMany(
    _value: Partial<BaseInsertion<FileEntity>>,
    _query: SearchQuery<FileEntity>
  ): Promise<{ count: number; data: FileEntity[] }> {
    throw new Error('Method not implemented.')
  }

  protected emitMutationEvent(entity: FileEntity): void {
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['files', FILE_ENTITY_KEY]
    })
  }
}
