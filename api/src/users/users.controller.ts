import { Body, Controller, Get, Patch, Post } from "@nestjs/common";
import { ApiBearerAuth, ApiTags } from "@nestjs/swagger";

import { CurrentUser } from "../auth/decorators/current-user.decorator";
import { UpdateUserDto } from "./dto/update-user.dto";
import { UserEntity } from "./entities/user.entity";
import { UsersService } from "./users.service";

@ApiBearerAuth()
@ApiTags("User")
@Controller("/user")
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  /**
   * Deactivate a user
   * @remarks This endpoint deactivates a user.
   * @throws {401} Unauthorized.
   * @throws {400} Bad Request.
   */
  @Post("deactivate")
  async deactivate(@CurrentUser() user: UserEntity) {
    return this.usersService.deactivate(user.id);
  }

  /**
   * Get a user
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Get()
  async findOne(@CurrentUser() user: UserEntity) {
    return user;
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
    return this.usersService.update(
      user.defaultOrgname,
      user.id,
      updateUserDto
    );
  }
}
