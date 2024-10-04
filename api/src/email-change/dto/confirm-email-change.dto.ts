import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class ConfirmEmailChangeDto {
  @ApiProperty({
    description:
      "The token used to change the e-mail. This should have been provided in an e-mail",
    example: "supersecre",
  })
  @IsString()
  token: string;
}
