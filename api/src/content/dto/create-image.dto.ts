import { PickType } from "@nestjs/swagger";

import { ImageEntity } from "../../content/entities/image.entity";

export class CreateImageDto extends PickType(ImageEntity, [
  "name",
  "width",
  "height",
  "prompt",
] as const) {}
