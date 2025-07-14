import type { WebsocketsService } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { ARTIFACT_ENTITY_KEY } from '@archesai/schemas'

import type { ArtifactRepository } from '#artifacts/artifact.repository'

export const createArtifactsService = (
  artifactRepository: ArtifactRepository,
  websocketsService: WebsocketsService
) =>
  createBaseService(
    artifactRepository,
    websocketsService,
    emitArtifactMutationEvent
  )

const emitArtifactMutationEvent = (
  entity: ArtifactEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, ARTIFACT_ENTITY_KEY]
  })
}

export type ArtifactsService = ReturnType<typeof createArtifactsService>

// async populateReadUrl(content: ArtifactEntity) {
//     const url = `https://storage.googleapis.com/archesai/storage/${content.organizationId}/`
//     if (!content.url?.startsWith(url)) {
//       this.logger.debug('url does not start with storage url', { content })
//       return content
//     }

//     const path = content.url.replace(url, '').split('?')[0] ?? ''
//     const [err, file] = await catchErrorAsync(
//       this.storageService.createSignedUrl(decodeURIComponent(path), 'read')
//     )

//     if (err) {
//       this.logger.warn('error getting signed url', { error: err })
//       content.url = ''
//       return content
//     }

//     this.logger.debug('got signed url', { file })
//     content.url = file.read ?? ''
//     return content
//   }

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
