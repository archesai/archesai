import { BadRequestException, Inject, Injectable } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { ContentRepository } from './content.repository'
import { ContentEntity, ContentModel } from './entities/content.entity'
import { ScraperService } from '../scraper/scraper.service'
import { v4 } from 'uuid'
import {
  IStorageService,
  STORAGE_SERVICE
} from '../storage/interfaces/storage-provider.interface'

@Injectable()
export class ContentService extends BaseService<
  ContentEntity,
  ContentModel,
  ContentRepository
> {
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: IStorageService,
    private contentRepository: ContentRepository,
    private websocketsService: WebsocketsService,
    private scraperService: ScraperService
  ) {
    super(contentRepository)
  }

  async create(
    data: Pick<ContentEntity, 'url' | 'text' | 'name' | 'orgname' | 'labels'>
  ) {
    let mimeType: string
    if (data.url) {
      mimeType = await this.scraperService.detectMimeType(data.url)
    } else if (data.text) {
      mimeType = 'text/plain'
    } else {
      throw new BadRequestException('Either url or text must be provided')
    }
    const previewBuffer = await this.scraperService.generateThumbnail(
      data.url,
      data.text,
      mimeType
    )
    const id = v4()
    const previewImage = await this.storageService.upload(
      data.orgname,
      `contents/${id}-preview.png`,
      {
        buffer: previewBuffer,
        mimetype: 'image/png',
        originalname: `${id}-preview.png`,
        size: previewBuffer.length
      } as Express.Multer.File
    )

    const content = await this.repository.create({
      name: data.name,
      url: data.url,
      text: data.text,
      orgname: data.orgname,
      mimeType,
      previewImage
    })
    const entity = this.toEntity(content)
    this.emitMutationEvent(entity)
    return content
  }

  async populateReadUrl(content: ContentModel) {
    const url = `https://storage.googleapis.com/archesai/storage/${content.orgname}/`
    if (content.url?.startsWith(url)) {
      const path = content.url.replace(url, '').split('?')[0]
      try {
        const read = await this.storageService.getSignedUrl(
          content.orgname,
          decodeURIComponent(path),
          'read'
        )
        content.url = read
      } catch (e) {
        this.logger.warn(e)
        content.url = ''
      }
    }
    return this.toEntity(content)
  }

  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ) {
    return this.contentRepository.query(orgname, embedding, topK, contentIds)
  }

  protected emitMutationEvent(entity: ContentEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'content']
    })
  }

  protected toEntity(model: ContentModel): ContentEntity {
    return new ContentEntity(model)
  }
}
