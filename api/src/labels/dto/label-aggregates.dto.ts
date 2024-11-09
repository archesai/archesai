import { ApiProperty } from "@nestjs/swagger";

import {
  GranularCount,
  GranularSum,
} from "../../common/dto/aggregated-field.dto";

export class LabelAggregates {
  @ApiProperty({
    description:
      "The number of credits used in chat labels over specific timeframes",
    required: false,
    type: [GranularSum],
  })
  credits: GranularSum[];

  @ApiProperty({
    description: "The number of labels created over specific timeframes",
    required: false,
    type: [GranularCount],
  })
  labelsCreated: GranularCount[];
}
