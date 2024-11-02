import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { TextChunk } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";

import { BaseEntity } from "../../common/base-entity.dto";

@Exclude()
export class TextChunkEntity extends BaseEntity implements TextChunk {
  @ApiHideProperty()
  contentId: string;

  @ApiHideProperty()
  orgname: string;

  @ApiProperty({
    description: "The job that created this vector record",
  })
  @Expose()
  text: string;

  constructor(textChunk: TextChunk) {
    super();
    Object.assign(this, textChunk);
  }
}
