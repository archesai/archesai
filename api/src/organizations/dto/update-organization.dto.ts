import { PickType } from "@nestjs/swagger";

import { CreateOrganizationDto } from "./create-organization.dto";

export class UpdateOrganizationDto extends PickType(CreateOrganizationDto, [
  "billingEmail",
] as const) {}
