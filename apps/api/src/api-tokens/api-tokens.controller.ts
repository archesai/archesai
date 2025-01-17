import { Body, Controller, Param, Post } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { BaseController } from '@/src/common/base.controller'
import { UserEntity } from '@/src/users/entities/user.entity'
import { ApiTokensService } from '@/src/api-tokens/api-tokens.service'
import { CreateApiTokenDto } from '@/src/api-tokens/dto/create-api-token.dto'
import { UpdateApiTokenDto } from '@/src/api-tokens/dto/update-api-token.dto'
import { ApiTokenEntity } from '@/src/api-tokens/entities/api-token.entity'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags(`API Tokens`)
@Authenticated()
@Controller('/organizations/:orgname/api-tokens')
export class ApiTokensController extends BaseController<
  ApiTokenEntity,
  CreateApiTokenDto,
  UpdateApiTokenDto,
  ApiTokensService
>(ApiTokenEntity, CreateApiTokenDto, UpdateApiTokenDto) {
  constructor(private readonly apiTokensService: ApiTokensService) {
    super(apiTokensService)
  }

  /**
   * Create a new API token
   * @remarks This endpoint requires the user to be authenticated
   */
  @Post()
  async create(
    @Param('orgname') orgname: string,
    @Body() createTokenDto: CreateApiTokenDto,
    @CurrentUser() currentUserDto: UserEntity
  ) {
    this.logger.debug(`creating ${this.itemType}`, {
      orgname,
      createTokenDto,
      currentUserDto
    })
    return this.apiTokensService.create({
      ...createTokenDto,
      username: currentUserDto.username,
      orgname
    })
  }
}
