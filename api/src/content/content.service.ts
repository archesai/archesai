import { BadRequestException, Inject, Injectable, Logger } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { STORAGE_SERVICE, StorageService } from '../storage/storage.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { ContentRepository } from './content.repository'
import { CreateContentDto } from './dto/create-content.dto'
import { UpdateContentDto } from './dto/update-content.dto'
import { ContentEntity, ContentModel } from './entities/content.entity'
import { ScraperService } from '../scraper/scraper.service'

@Injectable()
export class ContentService extends BaseService<
  ContentEntity,
  CreateContentDto,
  UpdateContentDto,
  ContentRepository,
  ContentModel
> {
  private logger = new Logger(ContentService.name)
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private contentRepository: ContentRepository,
    private websocketsService: WebsocketsService,
    private scraperService: ScraperService
  ) {
    super(contentRepository)
  }

  async create(orgname: string, createContentDto: CreateContentDto) {
    let mimeType: string
    if (createContentDto.url) {
      mimeType = await this.scraperService.detectMimeType(createContentDto.url)
    } else if (createContentDto.text) {
      mimeType = 'text/plain'
    } else {
      throw new BadRequestException('Either url or text must be provided')
    }
    let content = this.toEntity(
      await this.contentRepository.create(orgname, createContentDto, {
        mimeType
      })
    )
    const previewBuffer = await this.scraperService.generateThumbnail(
      content.url,
      content.text,
      content.mimeType
    )
    const url = await this.storageService.upload(
      orgname,
      `contents/${content.name}-preview.png`,
      {
        buffer: previewBuffer,
        mimetype: 'image/png',
        originalname: `${content.name}-preview.png`,
        size: previewBuffer.length
      } as Express.Multer.File
    )
    content = await this.setPreviewImage(orgname, content.id, url)
    this.emitMutationEvent(orgname)
    return content
  }

  async incrementCredits(orgname: string, id: string, credits: number) {
    return this.toEntity(
      await this.contentRepository.updateRaw(orgname, id, {
        credits: { increment: credits }
      })
    )
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

  async setPreviewImage(orgname: string, id: string, previewImage: string) {
    return this.toEntity(
      await this.contentRepository.updateRaw(orgname, id, { previewImage })
    )
  }

  async setTitle(orgname: string, id: string, title: string) {
    return this.toEntity(
      await this.contentRepository.updateRaw(orgname, id, {
        name: title
      })
    )
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit('update', {
      queryKey: ['organizations', orgname, 'content']
    })
  }

  protected toEntity(model: ContentModel): ContentEntity {
    return new ContentEntity(model)
  }
}
