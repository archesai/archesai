import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { VectorRecord } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";

import { BaseEntity } from "../../common/base-entity.dto";

@Exclude()
export class VectorRecordEntity extends BaseEntity implements VectorRecord {
  @ApiHideProperty()
  contentId: string;

  @ApiHideProperty()
  orgname: string;

  @ApiProperty({
    description: "The job that created this vector record",
  })
  @Expose()
  text: string;

  constructor(vectorRecord: VectorRecord) {
    super();
    Object.assign(this, vectorRecord);
  }
}
