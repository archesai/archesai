import { Body, Controller, Get, Post, Res } from "@nestjs/common";
import { UseGuards } from "@nestjs/common";
import { AuthGuard } from "@nestjs/passport";
import {
  ApiBadRequestResponse,
  ApiBearerAuth,
  ApiCreatedResponse,
  ApiExcludeEndpoint,
  ApiOperation,
  ApiResponse,
  ApiTags,
  ApiUnauthorizedResponse,
} from "@nestjs/swagger";
import { Response } from "express";

import { AuthService } from "./auth.service";
import {
  CurrentUser,
  CurrentUserDto,
} from "./decorators/current-user.decorator";
import { IsPublic } from "./decorators/is-public.decorator";
import { Roles } from "./decorators/roles.decorator";
import { ConfirmationTokenDto } from "./dto/confirmation-token.dto";
import { ConfirmationTokenWithNewPasswordDto } from "./dto/confirmation-token-with-new-password.dto";
import { EmailRequestDto } from "./dto/email-request.dto";
import { LoginDto } from "./dto/login.dto";
import { RegisterDto } from "./dto/register.dto";
import { TokenDto } from "./dto/token.dto";
import { EmailChangeService } from "./email-change.service";
import { EmailVerificationService } from "./email-verification.service";
import { LocalAuthGuard } from "./guards/local-auth.guard";
import { PasswordResetService } from "./password-reset.service";

@ApiTags("Authentication")
@Controller("auth")
export class AuthController {
  constructor(
    private readonly authService: AuthService,
    private passwordResetService: PasswordResetService,
    private emailVerificationService: EmailVerificationService,
    private emailChangeService: EmailChangeService
  ) {}

  @ApiOperation({
    description: "This endpoint will confirm your e-mail change with a token",
    summary: "Confirm e-mail change",
  })
  @ApiResponse({
    description: "E-mail change confirmed",
    status: 201,
    type: TokenDto,
  })
  @ApiBadRequestResponse()
  @IsPublic()
  @Post("/email-change/confirm")
  async emailChangeConfirm(
    @Body() confirmEmailChangeDto: ConfirmationTokenDto
  ) {
    return new TokenDto(
      await this.emailChangeService.confirm(confirmEmailChangeDto)
    );
  }

  @ApiBearerAuth()
  @ApiOperation({
    description: "This endpoint will request your e-mail change with a token",
    summary: "Confirm e-mail change",
  })
  @Post("/email-change/request")
  async emailChangeRequest(
    @CurrentUser() currentUserDto: CurrentUserDto,
    @Body() emailRequestDto: EmailRequestDto
  ) {
    return this.emailChangeService.request(currentUserDto.id, emailRequestDto);
  }

  @ApiOperation({
    description: "This endpoint will confirm your e-mail with a token",
    summary: "Confirm e-mail verification",
  })
  @ApiBadRequestResponse()
  @ApiUnauthorizedResponse()
  @ApiCreatedResponse({
    description: "E-mail verification confirmed",
    type: TokenDto,
  })
  @IsPublic()
  @Post("/email-verification/confirm")
  async emailVerificationConfirm(
    @Body() confirmEmailVerificationDto: ConfirmationTokenDto
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
  @Post("/email-verification/request")
  async emailVerificationRequest(@CurrentUser() user: CurrentUserDto) {
    return this.emailVerificationService.request(user.id);
  }

  @ApiOperation({ summary: "Login" })
  @ApiUnauthorizedResponse({ description: "Invalid credentials" })
  @ApiCreatedResponse({
    type: TokenDto,
  })
  @UseGuards(LocalAuthGuard)
  @IsPublic()
  @Post("/login")
  async login(
    @Body() loginDto: LoginDto,
    @CurrentUser() currentUserDto: CurrentUserDto,
    @Res({
      passthrough: true,
    })
    res: Response
  ): Promise<TokenDto> {
    const tokenDto = await this.authService.login(currentUserDto);
    await this.authService.setCookies(res, tokenDto);
    return tokenDto;
  }

  @ApiOperation({
    description: "Log out of the current session",
    summary: "Logout",
  })
  @IsPublic()
  @Post("/logout")
  async logout(
    @Res({
      passthrough: true,
    })
    res: Response
  ): Promise<void> {
    await this.authService.removeCookies(res);
  }

  @ApiOperation({
    description: "This endpoint will confirm your password change with a token",
    summary: "Confirm password change",
  })
  @ApiResponse({
    description: "Password change confirmed",
    status: 201,
    type: TokenDto,
  })
  @IsPublic()
  @Post("/password-reset/confirm")
  async passwordResetConfirm(
    @Body() confirmPasswordReset: ConfirmationTokenWithNewPasswordDto
  ): Promise<TokenDto> {
    return new TokenDto(
      await this.passwordResetService.confirm(confirmPasswordReset)
    );
  }

  @ApiOperation({
    description: "This endpoint will request a password reset link",
    summary: "Request password reset",
  })
  @IsPublic()
  @Post("/password-reset/request")
  async passwordResetRequest(@Body() emailRequestDto: EmailRequestDto) {
    await this.passwordResetService.request(emailRequestDto);
  }

  @ApiOperation({ summary: "Refresh Access Token" })
  @ApiResponse({
    description: "The new access token has been generated.",
    status: 201,
    type: TokenDto,
  })
  @ApiUnauthorizedResponse({ description: "Invalid refresh token" })
  @IsPublic()
  @Post("/refresh-token")
  async refreshToken(
    @Body("refreshToken") refreshToken: string
  ): Promise<TokenDto> {
    const tokens = await this.authService.refreshAccessToken(refreshToken);
    return tokens;
  }

  @ApiOperation({
    description:
      "This endpoint will register a new account and return a JWT token which should be provided in your auth headers",
    summary: "Register",
  })
  @ApiResponse({
    description: "User already exists with provided email",
    status: 409,
  })
  @ApiResponse({
    description: "User was created successfully",
    status: 201,
    type: TokenDto,
  })
  @IsPublic()
  @Post("/register")
  async register(@Body() registerDto: RegisterDto): Promise<TokenDto> {
    const user = await this.authService.register(registerDto);
    return this.authService.login(user);
  }

  @ApiExcludeEndpoint()
  @UseGuards(AuthGuard("firebase-auth"))
  @Post("firebase/callback")
  async zfirebaseAuthCallback(
    @CurrentUser() currentUserDto: CurrentUserDto
  ): Promise<TokenDto> {
    return this.authService.login(currentUserDto);
  }

  @ApiExcludeEndpoint()
  @UseGuards(AuthGuard("twitter"))
  @Get("twitter")
  async ztwitterAuth() {}

  @ApiExcludeEndpoint()
  @UseGuards(AuthGuard("twitter"))
  @Get("twitter/callback")
  async ztwitterAuthCallback(
    @CurrentUser() currentUserDto: CurrentUserDto
  ): Promise<TokenDto> {
    return this.authService.login(currentUserDto);
  }
}
