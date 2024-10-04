import { ApiProperty } from "@nestjs/swagger";
import { Exclude, Expose } from "class-transformer";

export enum Granularity {
  DAY = "day",
  MONTH = "month",
  WEEK = "week",
  YEAR = "year",
}

@Exclude()
export class GroupCount {
  @Expose()
  @ApiProperty()
  count: number;

  @Expose()
  @ApiProperty()
  value: string;
}

@Exclude()
export class GranularCount {
  @Expose()
  @ApiProperty()
  count: number;

  @Expose()
  @ApiProperty()
  from: Date;

  @Expose()
  @ApiProperty()
  to: Date;
}

@Exclude()
export class GranularSum {
  @Expose()
  @ApiProperty()
  from: Date;

  @Expose()
  @ApiProperty()
  sum: number;

  @Expose()
  @ApiProperty()
  to: Date;
}
