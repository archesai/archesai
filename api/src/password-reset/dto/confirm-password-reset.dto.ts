import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class ConfirmPasswordResetDto {
  @ApiProperty({
    description: "The new password",
    example: "newPassword",
  })
  @IsString()
  newPassword: string;

  @ApiProperty({
    description:
      "The token used to reset the password. This should have been provided in an e-mail",
    example: "supersecre",
  })
  @IsString()
  token: string;
}
