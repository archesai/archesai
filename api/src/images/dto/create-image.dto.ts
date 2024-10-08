import { PickType } from "@nestjs/swagger";

import { ImageEntity } from "../entities/image.entity";

export class CreateImageDto extends PickType(ImageEntity, [
  "name",
  "width",
  "height",
  "prompt",
] as const) {}
