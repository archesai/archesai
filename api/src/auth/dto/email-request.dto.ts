import { ApiProperty } from '@nestjs/swagger'
import { IsEmail, IsString } from 'class-validator'

export class EmailRequestDto {
  @ApiProperty({
    description: 'The e-mail to send the confirmation token to',
    example: 'user@archesai.com'
  })
  @IsEmail()
  @IsString()
  email: string
}
