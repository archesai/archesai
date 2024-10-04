// src/password-reset/password-reset.service.ts

import { Injectable } from "@nestjs/common";
import { ARTokenType } from "@prisma/client";
import * as bcrypt from "bcryptjs";

import { EmailService } from "../email/email.service";
import { getPasswordResetHtml } from "../email/templates";
import { PrismaService } from "../prisma/prisma.service";
import { ConfirmPasswordResetDto } from "./dto/confirm-password-reset.dto"; // HTML Email Templates
import { ARTokensService } from "../ar-tokens/ar-tokens.service";
import { AuthService } from "../auth/auth.service"; // Import TokenService
import { ConfigService } from "@nestjs/config";

import { RequestPasswordResetDto } from "./dto/request-password-reset.dto";

@Injectable()
export class PasswordResetService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly emailService: EmailService,
    private readonly authService: AuthService,
    private readonly arTokensService: ARTokensService,
    private readonly configService: ConfigService
  ) {}

  async confirm(confirmResetPasswordDto: ConfirmPasswordResetDto) {
    const { userId } = await this.arTokensService.verifyToken(
      ARTokenType.PASSWORD_RESET,
      confirmResetPasswordDto.token
    );

    const hashedPassword = await bcrypt.hash(
      confirmResetPasswordDto.newPassword,
      10
    );
    const user = await this.prisma.user.update({
      data: { password: hashedPassword },
      include: {
        authProviders: true,
        memberships: true,
      },
      where: { id: userId },
    });

    return this.authService.login(user);
  }

  async request(
    requestPasswordResetDto: RequestPasswordResetDto
  ): Promise<void> {
    const user = await this.prisma.user.findUnique({
      where: { email: requestPasswordResetDto.email },
    });
    if (!user) {
      return;
    }

    const token = await this.arTokensService.createToken(
      ARTokenType.PASSWORD_RESET,
      user.id,
      1
    );

    const resetLink = `${this.configService.get(
      "FRONTEND_HOST"
    )}/auth/confirm?type=password-reset&token=${token}`;

    const htmlContent = getPasswordResetHtml(resetLink);
    await this.emailService.sendMail({
      html: htmlContent,
      subject: "Password Reset Request",
      to: requestPasswordResetDto.email,
    });
  }
}
