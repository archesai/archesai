import { Content as _PrismaContent } from '@prisma/client'

import {
  _PrismaSubItemModel,
  SubItemEntity
} from '../../common/entities/base-sub-item.entity'
import { BaseEntity } from '../../common/entities/base.entity'
import { Expose } from 'class-transformer'
import { IsNumber, IsOptional, IsString, ValidateNested } from 'class-validator'

export type ContentModel = _PrismaContent & {
  children: _PrismaSubItemModel[]
  consumedBy: _PrismaSubItemModel[]
  labels: _PrismaSubItemModel[]
  parent: _PrismaSubItemModel | null
  producedBy: _PrismaSubItemModel | null
}

export class ContentEntity extends BaseEntity implements ContentModel {
  /**
   * The child content, if any
   */
  @ValidateNested({
    each: true
  })
  @Expose()
  children: SubItemEntity[]

  /**
   * The tool runs that consumed this content, if any
   */
  @ValidateNested({
    each: true
  })
  @Expose()
  consumedBy: SubItemEntity[]

  /**
   * The number of credits used to process this content
   * @example 0
   */
  @IsNumber()
  @Expose()
  credits: number

  /**
   * The content's description
   * @example 'my-file.pdf'
   */
  @IsOptional()
  @IsString()
  @Expose()
  description: null | string

  /**
   * The content's labels
   */
  @ValidateNested({
    each: true
  })
  @Expose()
  labels: SubItemEntity[]

  /**
   * The MIME type of the content
   * @example 'application/pdf'
   */
  @IsOptional()
  @IsString()
  @Expose()
  mimeType: null | string

  /**
   * The content's name
   * @example 'my-file.pdf'
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The organization name
   * @example 'my-organization'
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The parent content, if any
   */
  @IsOptional()
  @ValidateNested()
  @Expose()
  parent: null | SubItemEntity

  /**
   * The parent content ID, if this content is a child of another content
   * @example 'content-id'
   */
  @IsOptional()
  @IsString()
  @Expose()
  parentId: null | string

  /**
   * The preview image of the content
   * @example 'https://preview-image.com/example.png'
   */
  @IsOptional()
  @IsString()
  @Expose()
  previewImage: null | string

  /**
   * The toolRun that produced this content, if any
   */
  @IsOptional()
  @ValidateNested()
  @Expose()
  producedBy: null | SubItemEntity

  /**
   * The ID of the toolRun that produced this content, if any
   * @example 'toolRun-id'
   */
  @IsOptional()
  @IsString()
  @Expose()
  producedById: null | string

  /**
   * The content's text, if TEXT content
   * @example 'Hello world. I am a text.'
   */
  @IsOptional()
  @IsString()
  @Expose()
  text: null | string

  /**
   * The URL of the content, if AUDIO, VIDEO, IMAGE, or FILE content
   * @example 'https://example.com/example.mp4'
   */
  @IsOptional()
  @IsString()
  @Expose()
  url: null | string

  constructor(content: ContentModel) {
    super()
    Object.assign(this, content)
  }
}
