import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class ConfirmEmailVerificationDto {
  @ApiProperty({
    description:
      "The token used to verify the e-mail. This should have been provided in an e-mail",
    example: "supersecre",
  })
  @IsString()
  token: string;
}
