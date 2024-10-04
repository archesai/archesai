import { PickType } from "@nestjs/swagger";

import { ContentEntity } from "../entities/content.entity";

export class ContentFieldItem extends PickType(ContentEntity, ["id", "name"]) {}
