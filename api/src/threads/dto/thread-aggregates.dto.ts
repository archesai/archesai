import { ApiProperty } from "@nestjs/swagger";

import { GranularCount, GranularSum } from "../../common/aggregated-field.dto";

export class ThreadAggregates {
  @ApiProperty({
    description:
      "The number of credits used in chat threads over specific timeframes",
    required: false,
    type: [GranularSum],
  })
  credits: GranularSum[];

  @ApiProperty({
    description: "The number of threads created over specific timeframes",
    required: false,
    type: [GranularCount],
  })
  threadsCreated: GranularCount[];
}
