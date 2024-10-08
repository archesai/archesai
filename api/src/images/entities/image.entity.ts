import { ApiProperty } from "@nestjs/swagger";
import { Job } from "@prisma/client";
import { Content } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsBoolean, IsNumber, IsString } from "class-validator";

import { ContentEntity } from "../../content/entities/content.entity";
import { JobEntity } from "../../jobs/entities/job.entity";

export class ImageEntity extends ContentEntity implements Content {
  @ApiProperty({
    description: "The height of the image",
    example: 1024,
  })
  @Expose()
  @IsNumber()
  height: number;

  @ApiProperty({
    description: "The image prompt",
    example: "a person standing on the moon",
  })
  @Expose()
  @IsString()
  prompt: string;

  @ApiProperty({
    description: "Whether or not to use the init image",
    example: true,
  })
  @Expose()
  @IsBoolean()
  useInit: boolean;

  @ApiProperty({
    description: "The width of the image",
    example: 1024,
  })
  @Expose()
  @IsNumber()
  width: number;

  constructor(content: { job: Job } & Content) {
    super(content);
    Object.assign(this, content);
    this.job = new JobEntity(content.job);
  }
}
