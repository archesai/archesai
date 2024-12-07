import { HttpService } from '@nestjs/axios'
import { BadRequestException, Inject, Injectable, Logger } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { AxiosError } from 'axios'
import * as mime from 'mime-types'
import { catchError, firstValueFrom } from 'rxjs'

import { BaseService } from '../common/base.service'
import { STORAGE_SERVICE, StorageService } from '../storage/storage.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { ContentRepository } from './content.repository'
import { CreateContentDto } from './dto/create-content.dto'
import { UpdateContentDto } from './dto/update-content.dto'
import { ContentEntity, ContentModel } from './entities/content.entity'

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
    private configService: ConfigService,
    private httpService: HttpService
  ) {
    super(contentRepository)
  }

  async create(orgname: string, createContentDto: CreateContentDto) {
    let mimeType: string
    if (createContentDto.url) {
      mimeType = await this.detectMimeTypeFromUrl(createContentDto.url)
    } else {
      mimeType = 'text/plain'
    }

    // Create Content
    let content = this.toEntity(
      await this.contentRepository.create(orgname, createContentDto, {
        mimeType
      })
    )
    content = await this.getAndUploadPreview(orgname, content)
    this.emitMutationEvent(orgname)
    return content
  }

  async detectMimeTypeFromUrl(url: string) {
    try {
      // Extract the file name from the URL
      const urlObj = new URL(url)
      const pathname = urlObj.pathname
      const fileName = pathname.split('/').pop()

      if (!fileName) {
        throw new BadRequestException('Unable to extract file name from URL')
      }

      // Get MIME type based on file extension
      const mimeType = mime.lookup(fileName)
      return mimeType || null
    } catch (error) {
      this.logger.error('Failed to detect MIME type: ' + error.message)
      throw new BadRequestException('Failed to detect MIME type')
    }
  }

  async getAndUploadPreview(orgname: string, content: ContentEntity) {
    const { data } = await firstValueFrom(
      this.httpService
        .post(this.configService.get('LOADER_ENDPOINT') + '/getPreview', {
          text: content.text,
          url: content.url
        })
        .pipe(
          catchError((err: AxiosError) => {
            this.logger.error('Error hitting loader endpoint: ' + err.message)
            throw new BadRequestException()
          })
        )
    )
    const { preview } = data
    const previewFilename = `${content.name}-preview.png`
    const decodedImage = Buffer.from(preview, 'base64')
    const multerFile = {
      buffer: decodedImage,
      mimetype: 'image/png',
      originalname: previewFilename,
      size: decodedImage.length
    } as Express.Multer.File
    const url = await this.storageService.upload(orgname, `contents/${content.name}-preview.png`, multerFile)
    const updatedContent = await this.setPreviewImage(orgname, content.id, url)
    this.logger.log(`Upl image preview for ${content.name} at ${url}`)
    return updatedContent
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
        const read = await this.storageService.getSignedUrl(content.orgname, decodeURIComponent(path), 'read')
        content.url = read
      } catch (e) {
        this.logger.warn(e)
        content.url = ''
      }
    }
    return this.toEntity(content)
  }

  async query(orgname: string, embedding: number[], topK: number, contentIds?: string[]) {
    return this.contentRepository.query(orgname, embedding, topK, contentIds)
  }

  async setPreviewImage(orgname: string, id: string, previewImage: string) {
    return this.toEntity(await this.contentRepository.updateRaw(orgname, id, { previewImage }))
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
