import { Body, Controller, Param, Post } from '@nestjs/common'
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger'

import { CurrentUser } from '../auth/decorators/current-user.decorator'
import { BaseController } from '../common/base.controller'
import { UserEntity } from '../users/entities/user.entity'
import { ApiTokensService } from './api-tokens.service'
import { CreateApiTokenDto } from './dto/create-api-token.dto'
import { UpdateApiTokenDto } from './dto/update-api-token.dto'
import { ApiTokenEntity } from './entities/api-token.entity'

@ApiBearerAuth()
@ApiTags(`API Tokens`)
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
    return this.apiTokensService.create({
      ...createTokenDto,
      username: currentUserDto.username,
      orgname
    })
  }
}
