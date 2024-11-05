import { Body, Controller, Post } from "@nestjs/common";
import {
  ApiBadRequestResponse,
  ApiBearerAuth,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";
import { ApiOperation } from "@nestjs/swagger";

import {
  CurrentUser,
  CurrentUserDto,
} from "../auth/decorators/current-user.decorator";
import { IsPublic } from "../auth/decorators/is-public.decorator";
import { TokenDto } from "../auth/dto/token.dto";
import { ConfirmEmailChangeDto } from "./dto/confirm-email-change.dto";
import { RequestEmailChangeDto } from "./dto/request-email-change.dto";
import { EmailChangeService } from "./email-change.service";

@ApiTags("Authentication - Email Change")
@Controller("/auth/email-change")
export class EmailChangeController {
  constructor(private readonly emailChangeService: EmailChangeService) {}

  @ApiOperation({
    description: "This endpoint will confirm your e-mail change with a token",
    summary: "Confirm e-mail change",
  })
  @ApiResponse({
    description: "E-mail change confirmed",
    status: 201,
    type: TokenDto,
  })
  @ApiBadRequestResponse({
    description: "Invalid token",
    schema: {
      properties: {
        message: {
          example: "Invalid or expired token.",
          type: "string",
        },
        statusCode: {
          example: 400,
          type: "number",
        },
      },
    },
  })
  @IsPublic()
  @Post("confirm")
  async confirm(@Body() confirmEmailChangeDto: ConfirmEmailChangeDto) {
    return new TokenDto(
      await this.emailChangeService.confirm(confirmEmailChangeDto)
    );
  }

  @ApiBearerAuth()
  @ApiOperation({
    description: "This endpoint will request your e-mail change with a token",
    summary: "Confirm e-mail change",
  })
  @Post("request")
  async request(
    @CurrentUser() currentUserDto: CurrentUserDto,
    @Body() requestEmailChangeDto: RequestEmailChangeDto
  ) {
    return this.emailChangeService.request(
      currentUserDto.id,
      requestEmailChangeDto
    );
  }
}
