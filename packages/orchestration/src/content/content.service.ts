import type { WebsocketsService } from '@archesai/core'
import type { ContentEntity } from '@archesai/domain'
import type { StorageService } from '@archesai/storage'

import { BaseService, catchErrorAsync } from '@archesai/core'
import { CONTENT_ENTITY_KEY } from '@archesai/domain'

import type { ContentRepository } from '#content/content.repository'

/**
 * Service for content.
 */
export class ContentService extends BaseService<ContentEntity> {
  private readonly storageService: StorageService
  private readonly websocketsService: WebsocketsService

  constructor(
    contentRepository: ContentRepository,
    storageService: StorageService,
    websocketsService: WebsocketsService
  ) {
    super(contentRepository)
    this.storageService = storageService
    this.websocketsService = websocketsService
  }

  // override async create(data: ContentInsert): Promise<ContentEntity> {
  //   let mimeType: string
  //   if (data.url) {
  //     mimeType = await this.scraperService.detectMimeType(data.url)
  //   } else if (data.text) {
  //     mimeType = 'text/plain'
  //   } else {
  //     throw new BadRequestException('Either url or text must be provided')
  //   }
  //   const thumbnailBuffer = await this.scraperService.generateThumbnail(
  //     data.url,
  //     data.text,
  //     mimeType
  //   )
  //   const id = randomUUID()
  //   const thumbnailFile = await this.storageService.uploadFromFile(
  //     `contents/${id}-preview.png`,
  //     {
  //       buffer: thumbnailBuffer,
  //       mimetype: mimeType,
  //       originalname: `${id}-preview.png`
  //     }
  //   )

  //   const content = await this.contentRepository.create({
  //     mimeType,
  //     name: thumbnailFile.name,
  //     orgname: data.orgname,
  //     previewImage: thumbnailFile.read ?? '',
  //     text: data.text,
  //     url: data.url
  //   })
  //   const contentEntity = await this.findOne(content.id)
  //   this.emitMutationEvent(contentEntity)
  //   return contentEntity
  // }

  public async populateReadUrl(content: ContentEntity) {
    const url = `https://storage.googleapis.com/archesai/storage/${content.orgname}/`
    if (!content.url?.startsWith(url)) {
      this.logger.debug('url does not start with storage url', { content })
      return content
    }

    const path = content.url.replace(url, '').split('?')[0] ?? ''
    const [err, file] = await catchErrorAsync(
      this.storageService.createSignedUrl(decodeURIComponent(path), 'read')
    )

    if (err) {
      this.logger.warn('error getting signed url', { error: err })
      content.url = ''
      return content
    }

    this.logger.debug('got signed url', { file })
    content.url = file.read ?? ''
    return content
  }

  protected emitMutationEvent(entity: ContentEntity): void {
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname, CONTENT_ENTITY_KEY]
    })
  }
}
