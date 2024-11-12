import { Injectable, Logger, UnauthorizedException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtService } from "@nestjs/jwt";
import { AuthProviderType } from "@prisma/client";
import * as bcrypt from "bcryptjs";
import { Response } from "express";

import { UsersService } from "../users/users.service";
import { CurrentUserDto } from "./decorators/current-user.decorator";
import { RegisterDto } from "./dto/register.dto";
import { TokenDto } from "./dto/token.dto";

@Injectable()
export class AuthService {
  private readonly logger: Logger = new Logger(AuthService.name);

  constructor(
    protected jwtService: JwtService,
    protected usersService: UsersService,
    protected configService: ConfigService
  ) {}

  // Generate Access Token
  private generateAccessToken(userId: string) {
    this.logger.log("Generating access token for user: " + userId);
    return this.jwtService.sign(
      { sub: userId },
      {
        expiresIn: "15m", // Set access token expiration to 15 minutes
      }
    );
  }

  // Generate Refresh Token
  private generateRefreshToken(userId: string) {
    this.logger.log("Generating refresh token for user: " + userId);
    return this.jwtService.sign(
      { sub: userId },
      {
        expiresIn: "7d", // Set refresh token expiration to 7 days
        // secret: this.configService.get("JWT_REFRESH_SECRET"), // Use a different secret for refresh tokens
      }
    );
  }

  async login(user: CurrentUserDto) {
    this.logger.log("Logging in user: " + user.id);
    const accessToken = this.generateAccessToken(user.id);
    const refreshToken = this.generateRefreshToken(user.id);

    // Store refresh token in database
    await this.usersService.setRefreshToken(user.id, refreshToken);

    return {
      accessToken,
      refreshToken,
    };
  }

  // Refresh Access Token using Refresh Token
  async refreshAccessToken(refreshToken: string) {
    this.logger.log("Refreshing access token using refresh token");
    const payload = this.jwtService.verify(refreshToken, {
      // secret: this.configService.get("JWT_REFRESH_SECRET"),
    });

    const user = await this.usersService.findOne(null, payload.sub);

    if (!user || user.refreshToken !== refreshToken) {
      throw new UnauthorizedException();
    }

    // Generate new tokens
    const newAccessToken = this.generateAccessToken(user.id);
    const newRefreshToken = this.generateRefreshToken(user.id);

    // Update refresh token in the database
    await this.usersService.setRefreshToken(user.id, refreshToken);

    return {
      accessToken: newAccessToken,
      refreshToken: newRefreshToken, // Return new refresh token
    };
  }

  async register(registerDto: RegisterDto) {
    this.logger.log("Registering user: " + registerDto.email);
    const hashedPassword = await bcrypt.hash(registerDto.password, 10);
    const orgname =
      registerDto.email.split("@")[0] +
      "-" +
      Math.random().toString(36).substring(2, 6);
    const user = await this.usersService.create(null, {
      email: registerDto.email,
      emailVerified: this.configService.get("FEATURE_EMAIL") === false,
      password: hashedPassword,
      photoUrl: "",
      // username: registerDto.username,
      username: orgname,
    });
    return this.usersService.syncAuthProvider(
      user.email,
      AuthProviderType.LOCAL,
      user.email
    );
  }

  async removeCookies(res: Response) {
    res.clearCookie("archesai.accessToken");
    res.clearCookie("archesai.refreshToken");
  }

  async setCookies(res: Response, tokenDto: TokenDto) {
    res.cookie("archesai.accessToken", tokenDto.accessToken, {
      httpOnly: true,
      maxAge: 15 * 60 * 1000, // 15 minutes for access token
      sameSite: "strict",
      secure: this.configService.get("NODE_ENV") === "production",
    });
    res.cookie("archesai.refreshToken", tokenDto.refreshToken, {
      httpOnly: true,
      maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days for refresh token
      sameSite: "strict",
      secure: this.configService.get("NODE_ENV") === "production",
      signed: true,
    });
  }

  async verifyToken(token: string) {
    this.logger.log("Verifying jwt token: " + token);
    return this.jwtService.verify(token);
  }
}
