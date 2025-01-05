import { Expose } from 'class-transformer'
import { IsEmail } from 'class-validator'

export class EmailRequestDto {
  /**
   * The e-mail to send the confirmation token to
   * @example 'user@archesai.com'
   */
  @IsEmail()
  @Expose()
  email: string
}
