import { ApiHideProperty } from "@nestjs/swagger";
import { Exclude, Expose } from "class-transformer";

@Exclude()
export class BaseEntity {
  /**
   * The date that this item was created
   * @example '2023-07-11T21:09:20.895Z'
   */
  @Expose()
  createdAt: Date;

  /**
   * The ID of the item
   * @example 'item-id'
   */
  @Expose()
  id: string;

  @ApiHideProperty()
  updatedAt: Date;
}
