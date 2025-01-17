import { Body, Controller, Post } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { UserEntity } from '@/src/users/entities/user.entity'
import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { ConfirmationTokenDto } from '@/src/auth/dto/confirmation-token.dto'
import { EmailRequestDto } from '@/src/auth/dto/email-request.dto'
import { EmailChangeService } from '@/src/auth/services/email-change.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags('Authentication - Email Change')
@Controller('auth/email-change')
export class EmailChangeController {
  constructor(private emailChangeService: EmailChangeService) {}

  /**
   * Confirm e-mail change with a token
   * @remarks This endpoint will confirm your e-mail change with a token
   * @throws {400} BadRequestException
   */
  @Post('confirm')
  async emailChangeConfirm(
    @Body() confirmEmailChangeDto: ConfirmationTokenDto
  ): Promise<void> {
    await this.emailChangeService.confirm(confirmEmailChangeDto)
  }

  /**
   * Request e-mail change with a token
   * @remarks This endpoint will request your e-mail change with a token
   */
  @Authenticated()
  @Post('request')
  async emailChangeRequest(
    @CurrentUser() currentUserDto: UserEntity,
    @Body() emailRequestDto: EmailRequestDto
  ): Promise<void> {
    return this.emailChangeService.request(currentUserDto.id, emailRequestDto)
  }
}
