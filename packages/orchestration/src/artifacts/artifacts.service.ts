import type { WebsocketsService } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'
import type { StorageService } from '@archesai/storage'

import { BaseService, catchErrorAsync } from '@archesai/core'
import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'

import type { ArtifactRepository } from '#artifacts/artifact.repository'

/**
 * Service for content.
 */
export class ArtifactsService extends BaseService<ArtifactEntity> {
  private readonly storageService: StorageService
  private readonly websocketsService: WebsocketsService

  constructor(
    artifactRepository: ArtifactRepository,
    storageService: StorageService,
    websocketsService: WebsocketsService
  ) {
    super(artifactRepository)
    this.storageService = storageService
    this.websocketsService = websocketsService
  }

  // override async create(data: ContentInsert): Promise<ArtifactEntity> {
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

  //   const content = await this.artifactRepository.create({
  //     mimeType,
  //     name: thumbnailFile.name,
  //     organizationId: data.organizationId,
  //     previewImage: thumbnailFile.read ?? '',
  //     text: data.text,
  //     url: data.url
  //   })
  //   const contentEntity = await this.findOne(content.id)
  //   this.emitMutationEvent(contentEntity)
  //   return contentEntity
  // }

  public async populateReadUrl(content: ArtifactEntity) {
    const url = `https://storage.googleapis.com/archesai/storage/${content.organizationId}/`
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

  protected emitMutationEvent(entity: ArtifactEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, ARTIFACT_ENTITY_KEY]
    })
  }
}
