import { ApiProperty } from '@nestjs/swagger'
import { IsString } from 'class-validator'

import { ConfirmationTokenDto } from './confirmation-token.dto'

export class ConfirmationTokenWithNewPasswordDto extends ConfirmationTokenDto {
  @ApiProperty({
    description: 'The new password',
    example: 'newPassword'
  })
  @IsString()
  newPassword: string
}
