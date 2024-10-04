import { IntersectionType, PartialType, PickType } from "@nestjs/swagger";

import { DocumentEntity } from "../entities/document.entity";

export class CreateDocumentDto extends IntersectionType(
  PickType(DocumentEntity, ["name", "url"] as const),
  PartialType(PickType(DocumentEntity, ["delimiter", "chunkSize"] as const))
) {
  chunkSize?: number = 200;
  delimiter?: string = "";
}
