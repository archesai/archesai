import { PickType } from "@nestjs/swagger";

import { OrganizationEntity } from "../entities/organization.entity";
export class CreateOrganizationDto extends PickType(OrganizationEntity, [
  "orgname",
  "billingEmail",
]) {}
