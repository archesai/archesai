// src/email-change/email-change.service.ts

import {
  BadRequestException,
  ConflictException,
  Injectable,
} from "@nestjs/common";
import { ARTokenType } from "@prisma/client";

import { EmailService } from "../email/email.service";
import { getEmailChangeConfirmationHtml } from "../email/templates";
import { RequestEmailChangeDto } from "./dto/request-email-change.dto"; // HTML Email Templates
import { ConfigService } from "@nestjs/config";

import { ARTokensService } from "../ar-tokens/ar-tokens.service";
import { UsersService } from "../users/users.service"; // Import TokenService
import { AuthService } from "../auth/auth.service";
import { ConfirmEmailChangeDto } from "./dto/confirm-email-change.dto";

@Injectable()
export class EmailChangeService {
  constructor(
    private readonly emailService: EmailService,
    private readonly usersService: UsersService,
    private readonly configService: ConfigService,
    private readonly arTokensService: ARTokensService,
    private readonly authService: AuthService
  ) {}

  async confirm(confirmEmailChangeDto: ConfirmEmailChangeDto) {
    const { additionalData, userId } = await this.arTokensService.verifyToken(
      ARTokenType.EMAIL_CHANGE,
      confirmEmailChangeDto.token
    );

    const newEmail = additionalData?.newEmail;
    if (!newEmail) {
      throw new BadRequestException("New email is missing.");
    }

    const user = await this.usersService.updateEmail(userId, newEmail);
    return this.authService.login(user);
  }

  async request(
    userId: string,
    requestEmailChangeDto: RequestEmailChangeDto
  ): Promise<void> {
    const user = await this.usersService.findOne(null, userId);
    let newEmailInUse = false;
    try {
      await this.usersService.findOneByEmail(requestEmailChangeDto.email);
      newEmailInUse = true;
    } catch (e) {}
    if (newEmailInUse) {
      throw new ConflictException("New email is already in use.");
    }

    // Generate an email change token (expires in 24 hours) with additional data
    const token = await this.arTokensService.createToken(
      ARTokenType.EMAIL_CHANGE,
      user.id,
      24, // 24 hours expiry
      { newEmail: requestEmailChangeDto.email }
    );

    // Create an email change confirmation link containing the token
    const changeEmailLink = `${this.configService.get(
      "FRONTEND_HOST"
    )}/auth/confirm?type=email-change&token=${token}`;

    // Generate the HTML content for the email
    const htmlContent = getEmailChangeConfirmationHtml(
      changeEmailLink,
      user.email
    );

    // Send confirmation email to the new email address
    await this.emailService.sendMail({
      html: htmlContent,
      subject: "Confirm Your Email Change",
      to: requestEmailChangeDto.email,
    });
  }
}
