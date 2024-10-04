import { Injectable, UnauthorizedException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtService } from "@nestjs/jwt";
import { AuthProviderType } from "@prisma/client";
import * as bcrypt from "bcryptjs";

import { UsersService } from "../users/users.service";
import { CurrentUserDto } from "./decorators/current-user.decorator";
import { RegisterDto } from "./dto/register.dto";

@Injectable()
export class AuthService {
  constructor(
    protected jwtService: JwtService,
    protected usersService: UsersService,
    protected configService: ConfigService
  ) {}

  // Generate Access Token
  private generateAccessToken(userId: string) {
    return this.jwtService.sign(
      { sub: userId },
      {
        expiresIn: "15m", // Set access token expiration to 15 minutes
      }
    );
  }

  // Generate Refresh Token
  private generateRefreshToken(userId: string) {
    return this.jwtService.sign(
      { sub: userId },
      {
        expiresIn: "7d", // Set refresh token expiration to 7 days
        // secret: this.configService.get("JWT_REFRESH_SECRET"), // Use a different secret for refresh tokens
      }
    );
  }

  async login(user: CurrentUserDto) {
    const accessToken = this.generateAccessToken(user.id);
    const refreshToken = this.generateRefreshToken(user.id);

    // Store refresh token in database
    await this.usersService.updateRefreshToken(user.id, refreshToken);

    return {
      accessToken,
      refreshToken,
    };
  }

  // Refresh Access Token using Refresh Token
  async refreshAccessToken(refreshToken: string) {
    const payload = this.jwtService.verify(refreshToken, {
      // secret: this.configService.get("JWT_REFRESH_SECRET"),
    });

    const user = await this.usersService.findOne(payload.sub);

    if (!user || user.refreshToken !== refreshToken) {
      throw new UnauthorizedException();
    }

    // Generate new tokens
    const newAccessToken = this.generateAccessToken(user.id);
    const newRefreshToken = this.generateRefreshToken(user.id);

    // Update refresh token in the database
    await this.usersService.updateRefreshToken(user.id, newRefreshToken);

    return {
      accessToken: newAccessToken,
      refreshToken: newRefreshToken, // Return new refresh token
    };
  }

  async register(registerDto: RegisterDto) {
    const hashedPassword = await bcrypt.hash(registerDto.password, 10);
    console.log("FEATURE_EMAIL", this.configService.get("FEATURE_EMAIL"));
    const user = await this.usersService.create({
      email: registerDto.email,
      emailVerified: this.configService.get("FEATURE_EMAIL") === false,
      password: hashedPassword,
      photoUrl: "",
      username: registerDto.username,
    });
    return this.usersService.syncAuthProvider(
      user.email,
      AuthProviderType.LOCAL,
      user.email
    );
  }

  async verifyToken(token: string) {
    return this.jwtService.verify(token);
  }
}
