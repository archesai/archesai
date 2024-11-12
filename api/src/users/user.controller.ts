import { Body, Controller, Get, Patch, Post } from "@nestjs/common";
import {
  ApiBearerAuth,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";

import { CurrentUser } from "../auth/decorators/current-user.decorator";
import {
  ApiCrudOperation,
  Operation,
} from "../common/decorators/api-crud-operation.decorator";
import { UpdateUserDto } from "./dto/update-user.dto";
import { UserEntity } from "./entities/user.entity";
import { UsersService } from "./users.service";

@ApiBearerAuth()
@ApiTags("User")
@Controller("user")
export class UserController {
  constructor(private readonly usersService: UsersService) {}

  @ApiResponse({
    description: "Unauthorized",
    status: 401,
  })
  @ApiResponse({
    description: "User was deleted successfully",
    status: 201,
  })
  @ApiOperation({ summary: "Deactivate" })
  @Post("/deactivate")
  async deactivate(@CurrentUser() user: UserEntity) {
    return this.usersService.deactivate(user.id);
  }

  @ApiCrudOperation(Operation.GET, "user", UserEntity, true)
  @Get()
  async findOne(@CurrentUser() user: UserEntity) {
    return user;
  }

  @ApiCrudOperation(Operation.UPDATE, "user", UserEntity, true)
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
