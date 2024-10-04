import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Content, ContentType, Job, Organization } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsEnum, IsNumber, IsObject, IsString } from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto";
import { JobEntity } from "../../jobs/entities/job.entity";

@Exclude()
export class ContentEntity extends BaseEntity implements Content {
  @ApiProperty({
    description: "The content's annotations",
    example: { key: "value" },
  })
  @Expose()
  @IsObject()
  annotations: any;

  @ApiProperty()
  @Expose()
  @IsObject()
  buildArgs: any;

  @ApiProperty({
    description: "The number of credits used to process this content",
    example: 0,
  })
  @Expose()
  @IsNumber()
  credits: number;

  @ApiProperty({
    description: "The animation's name",
    example: "my-file.pdf",
  })
  @Expose()
  @IsString()
  description: string;

  @ApiProperty({
    description: "This job associated with this content's build process",
    type: JobEntity,
  })
  @Expose()
  job: JobEntity;

  @ApiHideProperty()
  jobId: string;

  @ApiProperty({ example: "application/pdf" })
  @Expose()
  mimeType: string;

  @ApiProperty({
    description: "The animation's name",
    example: "my-file.pdf",
  })
  @Expose()
  @IsString()
  name: string;

  // Private Properties
  @ApiHideProperty()
  organization: Organization;

  @ApiProperty({
    description: "The organization name",
    example: "my-organization",
  })
  @Expose()
  @IsString()
  orgname: string;

  @ApiProperty({
    description: "The preview image of the animation",
    example: "https://preview-image.com/example.png",
  })
  @Expose()
  @IsString()
  previewImage: string;

  @ApiProperty({
    description: "The content's text",
    example: "Hello, world!",
  })
  text: string;

  @ApiProperty({ enum: ContentType, example: "video/mp4" })
  @Expose()
  @IsEnum(ContentType)
  type: ContentType;

  @ApiProperty({ example: "https://example.com/example.mp4" })
  @Expose()
  @IsString()
  url: string;

  constructor(content: { job: Job } & Content) {
    super();
    Object.assign(this, content);
    this.job = new JobEntity(content.job);
  }
}
