import { Body, Controller, Get, Post } from "@nestjs/common";
import { UseGuards } from "@nestjs/common";
import { AuthGuard } from "@nestjs/passport";
import {
  ApiCreatedResponse,
  ApiExcludeEndpoint,
  ApiOperation,
  ApiResponse,
  ApiTags,
  ApiUnauthorizedResponse,
} from "@nestjs/swagger";

import { AuthService } from "./auth.service";
import {
  CurrentUser,
  CurrentUserDto,
} from "./decorators/current-user.decorator";
import { IsPublic } from "./decorators/is-public.decorator";
import { LoginDto } from "./dto/login.dto";
import { RegisterDto } from "./dto/register.dto";
import { TokenDto } from "./dto/token.dto";

@IsPublic()
@ApiTags("Authentication")
@Controller("auth")
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @ApiOperation({ summary: "Login" })
  @ApiUnauthorizedResponse({ description: "Invalid credentials" })
  @ApiCreatedResponse({
    type: TokenDto,
  })
  @UseGuards(AuthGuard("local"))
  @Post("/login")
  async login(
    @Body() loginDto: LoginDto,
    @CurrentUser() currentUserDto: CurrentUserDto
  ): Promise<TokenDto> {
    return this.authService.login(currentUserDto);
  }

  @ApiOperation({ summary: "Refresh Access Token" })
  @ApiResponse({
    description: "The new access token has been generated.",
    status: 201,
    type: TokenDto,
  })
  @ApiUnauthorizedResponse({ description: "Invalid refresh token" })
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
    status: 404,
  })
  @ApiResponse({
    description: "User was created successfully",
    status: 201,
    type: TokenDto,
  })
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
