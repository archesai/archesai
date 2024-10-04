import { Body, Controller, Post } from "@nestjs/common";
import {
  ApiBearerAuth,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";

import { CurrentUser } from "../auth/decorators/current-user.decorator";
import { CurrentUserDto } from "../auth/decorators/current-user.decorator";
import { IsPublic } from "../auth/decorators/is-public.decorator";
import { Roles } from "../auth/decorators/roles.decorator";
import { TokenDto } from "../auth/dto/token.dto";
import { ConfirmEmailVerificationDto } from "./dto/confirm-password-reset.dto";
import { EmailVerificationService } from "./email-verification.service";

@ApiTags("Authentication - Email Verification")
@Controller("/auth/email-verification")
export class EmailVerificationController {
  constructor(
    private readonly emailVerificationService: EmailVerificationService
  ) {}

  @ApiOperation({
    description: "This endpoint will confirm your e-mail with a token",
    summary: "Confirm e-mail verification",
  })
  @ApiResponse({ description: "Already Verified", status: 400 })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({
    description: "E-mail verification confirmed",
    status: 201,
    type: TokenDto,
  })
  @IsPublic()
  @Post("confirm")
  async confirm(
    @Body() confirmEmailVerificationDto: ConfirmEmailVerificationDto
  ) {
    return new TokenDto(
      await this.emailVerificationService.confirm(confirmEmailVerificationDto)
    );
  }

  @ApiBearerAuth()
  @ApiOperation({
    description:
      "This endpoint will send an e-mail verification link to you. ADMIN ONLY.",
    summary: "Resend e-mail verification",
  })
  @ApiResponse({ description: "Already Verified", status: 400 })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({
    description: "E-mail verification link sent",
    status: 201,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Roles("ADMIN")
  @Post("request")
  async request(@CurrentUser() user: CurrentUserDto) {
    return this.emailVerificationService.request(user.id);
  }
}
