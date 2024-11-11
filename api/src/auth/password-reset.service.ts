// src/password-reset/password-reset.service.ts

import { Injectable } from "@nestjs/common";
import { ARTokenType } from "@prisma/client";
import * as bcrypt from "bcryptjs";

import { EmailService } from "../email/email.service";
import { getPasswordResetHtml } from "../email/templates";
import { PrismaService } from "../prisma/prisma.service";
import { ARTokensService } from "./ar-tokens.service";
import { AuthService } from "./auth.service"; // Import TokenService
import { ConfigService } from "@nestjs/config";

import { ConfirmationTokenWithNewPasswordDto } from "./dto/confirmation-token-with-new-password.dto";
import { EmailRequestDto } from "./dto/email-request.dto";

@Injectable()
export class PasswordResetService {
  constructor(
    private readonly prisma: PrismaService,
    private readonly emailService: EmailService,
    private readonly authService: AuthService,
    private readonly arTokensService: ARTokensService,
    private readonly configService: ConfigService
  ) {}

  async confirm(
    confirmationTokenWithNewPasswordDto: ConfirmationTokenWithNewPasswordDto
  ) {
    const { userId } = await this.arTokensService.verifyToken(
      ARTokenType.PASSWORD_RESET,
      confirmationTokenWithNewPasswordDto.token
    );

    const hashedPassword = await bcrypt.hash(
      confirmationTokenWithNewPasswordDto.newPassword,
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

  async request(emailRequestDto: EmailRequestDto): Promise<void> {
    const user = await this.prisma.user.findUnique({
      where: { email: emailRequestDto.email },
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
      to: emailRequestDto.email,
    });
  }
}
