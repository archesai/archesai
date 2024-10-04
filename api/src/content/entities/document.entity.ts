import { ApiProperty } from "@nestjs/swagger";
import { Content, Job } from "@prisma/client";
import { Expose } from "class-transformer";
import { IsNumber, IsOptional, IsString } from "class-validator";

import { ContentEntity } from "../../content/entities/content.entity";
import { JobEntity } from "../../jobs/entities/job.entity";

export class DocumentEntity extends ContentEntity {
  @ApiProperty({
    default: 200,
    description: "The size of the documents text segments",
    example: 200,
    required: false,
  })
  @Expose()
  @IsOptional()
  @IsNumber()
  chunkSize: number;

  @ApiProperty({
    default: "",
    description:
      "The delimiter used to separate the document into text segments. If left blank, only chunkSize will be used.",
    example: "",
    required: false,
  })
  @Expose()
  @IsString()
  @IsOptional()
  delimiter: string;

  @ApiProperty({ example: "ajskdflasb==" })
  @Expose()
  @IsString()
  md5Hash: string;

  @ApiProperty({
    description: "The document size in bytes",
    example: "1125162",
  })
  @Expose()
  @IsString()
  size: string;

  @ApiProperty({ deprecated: true, example: "THIS FIELD IS DEPRECATED" })
  @Expose()
  @IsString()
  sourceUrl: string;

  @ApiProperty({
    description: "The document's summary",
    example: "This is a summary of your file...",
  })
  @Expose()
  @IsString()
  summary: string;

  constructor(content: { job: Job } & Content) {
    super(content);
    Object.assign(this, content);
    this.job = new JobEntity(content.job);
  }
}
