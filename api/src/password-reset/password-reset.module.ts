import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { ARTokensModule } from "../ar-tokens/ar-tokens.module";
import { AuthModule } from "../auth/auth.module";
import { EmailModule } from "../email/email.module";
import { PrismaModule } from "../prisma/prisma.module";
import { PasswordResetController } from "./password-reset.controller";
import { PasswordResetService } from "./password-reset.service";

@Module({
  controllers: [PasswordResetController],
  imports: [
    EmailModule,
    PrismaModule,
    AuthModule,
    ARTokensModule,
    ConfigModule,
  ],
  providers: [PasswordResetService],
})
export class PasswordResetModule {}
