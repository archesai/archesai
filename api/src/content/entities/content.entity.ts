import { ApiProperty } from "@nestjs/swagger";
import { Content } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsNumber, IsString } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";

@Exclude()
export class ContentEntity extends BaseEntity implements Content {
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
    example: "https://example.com/example.mp4",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  url: null | string;

  constructor(content: Content) {
    super();
    Object.assign(this, content);
  }
}
