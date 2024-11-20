import { ApiProperty } from "@nestjs/swagger";
import { Exclude, Expose } from "class-transformer";
import { IsString } from "class-validator";

export type _PrismaSubItemModel = {
  id: string;
  name: string;
};

@Exclude()
export class SubItemEntity implements _PrismaSubItemModel {
  @ApiProperty({
    description: "The id of the item",
    example: "item-id",
  })
  @Expose()
  @IsString()
  id: string;

  @ApiProperty({
    description: "The name of the item",
    example: "item-name",
  })
  @Expose()
  @IsString()
  name: string;

  constructor(subItem: _PrismaSubItemModel) {
    Object.assign(this, subItem);
  }
}
