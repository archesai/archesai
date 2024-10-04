import { ApiProperty } from "@nestjs/swagger";
import { IsString } from "class-validator";

export class RequestEmailChangeDto {
  @ApiProperty({
    description: "The e-mail address to update to",
    example: "user@archesai.com",
  })
  @IsString()
  email: string;
}
