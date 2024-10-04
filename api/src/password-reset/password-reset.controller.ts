// src/password-reset/password-reset.controller.ts

import { Body, Controller, Post } from "@nestjs/common";
import { ApiOperation, ApiResponse, ApiTags } from "@nestjs/swagger";

import { IsPublic } from "../auth/decorators/is-public.decorator";
import { TokenDto } from "../auth/dto/token.dto";
import { ConfirmPasswordResetDto } from "./dto/confirm-password-reset.dto";
import { RequestPasswordResetDto } from "./dto/request-password-reset.dto";
import { PasswordResetService } from "./password-reset.service";

@ApiTags("Authentication - Password Reset")
@IsPublic()
@Controller("/auth/password-reset")
export class PasswordResetController {
  constructor(private readonly passwordResetService: PasswordResetService) {}

  @ApiOperation({
    description: "This endpoint will confirm your password change with a token",
    summary: "Confirm password change",
  })
  @ApiResponse({
    description: "Password change confirmed",
    status: 201,
    type: TokenDto,
  })
  @Post("confirm")
  async confirm(
    @Body() confirmResetPasswordDto: ConfirmPasswordResetDto
  ): Promise<TokenDto> {
    return new TokenDto(
      await this.passwordResetService.confirm(confirmResetPasswordDto)
    );
  }

  @ApiOperation({
    description: "This endpoint will request a password reset link",
    summary: "Request password reset",
  })
  @Post("request")
  async request(@Body() requestPasswordResetDto: RequestPasswordResetDto) {
    await this.passwordResetService.request(requestPasswordResetDto);
  }
}
