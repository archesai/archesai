import { Body, Controller, Get, Patch, Post } from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { CurrentUser } from '@/src/auth/decorators/current-user.decorator'
import { UpdateUserDto } from '@/src/users/dto/update-user.dto'
import { UserEntity } from '@/src/users/entities/user.entity'
import { UsersService } from '@/src/users/users.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'

@ApiTags('User')
@Authenticated()
@Controller('/user')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  /**
   * Deactivate a user
   * @remarks This endpoint deactivates a user.
   * @throws {401} Unauthorized.
   * @throws {400} Bad Request.
   */
  @Post('deactivate')
  async deactivate(@CurrentUser() user: UserEntity) {
    return this.usersService.deactivate(user.id)
  }

  /**
   * Get a user
   * @remarks This endpoint returns a user.
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Get()
  async findOne(@CurrentUser() user: UserEntity) {
    return user
  }

  /**
   * Update a user
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Patch()
  async update(
    @CurrentUser() user: UserEntity,
    @Body() updateUserDto: UpdateUserDto
  ) {
    return this.usersService.update(user.id, updateUserDto)
  }
}
