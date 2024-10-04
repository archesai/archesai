import { Body, Controller, Get, Patch, Post } from "@nestjs/common";
import {
  ApiBearerAuth,
  ApiNotFoundResponse,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";

import {
  CurrentUser,
  CurrentUserDto,
} from "../auth/decorators/current-user.decorator";
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
  async deactivate(@CurrentUser() user: CurrentUserDto) {
    return this.usersService.deactivate(user.id);
  }

  @ApiOperation({
    description:
      "This endpoint can be used to find out about the currently authorized user. USER and ADMIN can use endpoint.",
    summary: "Get current user ",
  })
  @ApiNotFoundResponse()
  @ApiResponse({
    description: "User was successfully returned",
    status: 200,
    type: UserEntity,
  })
  @Get()
  async findOne(@CurrentUser() user: CurrentUserDto) {
    return new UserEntity(user);
  }

  @ApiOperation({
    description:
      "This endpoint can be used to update the currently authorized user. ADMIN ONLY.",
    summary: "Update current user",
  })
  @ApiNotFoundResponse()
  @ApiResponse({
    description: "User was successfully updated",
    status: 200,
    type: UserEntity,
  })
  @Patch()
  async update(
    @CurrentUser() user: CurrentUserDto,
    @Body() updateUserDto: UpdateUserDto
  ) {
    return new UserEntity(
      await this.usersService.update(user.id, updateUserDto)
    );
  }
}
