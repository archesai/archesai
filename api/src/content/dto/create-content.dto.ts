import { PickType } from "@nestjs/swagger";

import { ContentEntity } from "../entities/content.entity";

export class CreateContentDto extends PickType(ContentEntity, [
  "name",
  "url",
] as const) {}
