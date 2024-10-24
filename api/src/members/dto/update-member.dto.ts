import { PartialType, PickType } from "@nestjs/swagger";

import { CreateMemberDto } from "./create-member.dto";

export class UpdateMemberDto extends PartialType(
  PickType(CreateMemberDto, ["role", "inviteEmail"] as const)
) {}
