import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { ARTokensModule } from "../ar-tokens/ar-tokens.module";
import { AuthModule } from "../auth/auth.module";
import { EmailModule } from "../email/email.module";
import { UsersModule } from "../users/users.module";
import { EmailVerificationController } from "./email-verification.controller";
import { EmailVerificationService } from "./email-verification.service";

@Module({
  controllers: [EmailVerificationController],
  imports: [EmailModule, UsersModule, ConfigModule, AuthModule, ARTokensModule],
  providers: [EmailVerificationService],
})
export class EmailVerificationModule {}
