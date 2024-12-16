import { IsString } from 'class-validator'

import { ConfirmationTokenDto } from './confirmation-token.dto'
import { Expose } from 'class-transformer'

export class ConfirmationTokenWithNewPasswordDto extends ConfirmationTokenDto {
  /**
   * The new password
   * @example 'newPassword'
   */
  @IsString()
  @Expose()
  newPassword: string
}
