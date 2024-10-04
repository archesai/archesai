import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class RequestPasswordResetDto {
  @ApiProperty({
    description: "The e-mail address to send the password reset link to",
    example: "user@archesai.com",
  })
  @IsString()
  email: string;
}
