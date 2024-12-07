import { IsString } from 'class-validator'

import { ConfirmationTokenDto } from './confirmation-token.dto'

export class ConfirmationTokenWithNewPasswordDto extends ConfirmationTokenDto {
  /**
   * The new password
   * @example 'newPassword'
   */
  @IsString()
  newPassword: string
}
