import { ApiProperty } from "@nestjs/swagger";
import { Content as _PrismaContent } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsNumber, IsString } from "class-validator";

import { BaseEntity } from "../../common/dto/base.entity.dto";

export type ContentModel = _PrismaContent;

@Exclude()
export class ContentEntity extends BaseEntity implements ContentModel {
  @ApiProperty({
    description: "The number of credits used to process this content",
    example: 0,
  })
  @Expose()
  @IsNumber()
  credits: number;

  @ApiProperty({
    description: "The content's description",
    example: "my-file.pdf",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  description: null | string;

  @ApiProperty({
    description: "The MIME type of the content",
    example: "application/pdf",
    required: false,
    type: String,
  })
  @Expose()
  mimeType: null | string;

  @ApiProperty({
    description: "The content's name",
    example: "my-file.pdf",
  })
  @Expose()
  @IsString()
  name: string;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  @IsString()
  orgname: string;

  @ApiProperty({
    description: "The preview image of the content",
    example: "https://preview-image.com/example.png",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  previewImage: null | string;

  @ApiProperty({
    description: "The content's text, if TEXT content",
    example: "Hello world. I am a text.",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  text: null | string;

  @ApiProperty({
    description:
      "The URL of the content, if AUDIO, VIDEO, IMAGE, or FILE content",
    example: "https://example.com/example.mp4",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  url: null | string;

  constructor(content: ContentModel) {
    super();
    Object.assign(this, content);
  }
}
