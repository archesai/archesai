import { HttpModule } from "@nestjs/axios";
import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { JwtModule } from "@nestjs/jwt";

import { ApiTokensModule } from "../api-tokens/api-tokens.module";
import { ARTokensModule } from "../ar-tokens/ar-tokens.module";
import { EmailModule } from "../email/email.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { UsersModule } from "../users/users.module";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";
import { EmailChangeService } from "./email-change.service";
import { EmailVerificationService } from "./email-verification.service";
import { PasswordResetService } from "./password-reset.service";
import { ApiKeyStrategy } from "./strategies/api-key.strategy";
import { FirebaseStrategy } from "./strategies/firebase.strategy";
import { JwtStrategy } from "./strategies/jwt.strategy";
import { LocalStrategy } from "./strategies/local-strategy";
import { TwitterStrategy } from "./strategies/twitter.strategy";

@Module({
  controllers: [AuthController],
  exports: [AuthService],
  imports: [
    HttpModule,
    ConfigModule,
    UsersModule,
    forwardRef(() => ApiTokensModule),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        secret: configService.get<string>("JWT_API_TOKEN_SECRET"),
      }),
    }),
    OrganizationsModule,
    PrismaModule,
    EmailModule,
    ARTokensModule,
  ],
  providers: [
    // Services
    AuthService,
    // Strategies
    LocalStrategy,
    JwtStrategy,
    FirebaseStrategy,
    TwitterStrategy,
    ApiKeyStrategy,
    // Additional Services
    PasswordResetService,
    EmailVerificationService,
    EmailChangeService,
  ],
})
export class AuthModule {}
