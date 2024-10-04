import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Exclude, Expose } from "class-transformer";

export class BaseEntity {
  @Expose()
  @ApiProperty({
    description: "The creation date of this item",
    example: "2023-07-11T21:09:20.895Z",
  })
  createdAt: Date;

  @Expose()
  @ApiProperty({
    description: "The item's unique identifier",
    example: "32411590-a8e0-11ed-afa1-0242ac120002",
  })
  id: string;

  @Exclude()
  @ApiHideProperty()
  updatedAt: Date;
}
