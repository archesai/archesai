import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Citation, Message } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";

@Exclude()
export class CitationEntity implements Citation {
  @ApiHideProperty()
  content: Message;

  @ApiHideProperty()
  contentId: string;

  @ApiHideProperty()
  id: string;

  @ApiHideProperty()
  message: Message;

  @ApiHideProperty()
  messageId: string;

  // Public Properties
  @ApiProperty({
    description: "The similarity of this source to the query",
    example: 0.82,
  })
  @Expose()
  similarity: number;

  constructor(
    citation: {
      message: Message;
    } & Citation
  ) {
    Object.assign(this, citation);
  }
}
