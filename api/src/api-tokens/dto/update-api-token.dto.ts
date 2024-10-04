import { PartialType } from "@nestjs/swagger";

import { CreateApiTokenDto } from "./create-api-token.dto";

export class UpdateApiTokenDto extends PartialType(CreateApiTokenDto) {}
