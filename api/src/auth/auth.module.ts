import { HttpModule } from "@nestjs/axios";
import { forwardRef, Module } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtModule } from "@nestjs/jwt";
import { PassportModule } from "@nestjs/passport";

import { ApiTokensModule } from "../api-tokens/api-tokens.module";
import { EmailModule } from "../email/email.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { PrismaModule } from "../prisma/prisma.module";
import { UsersModule } from "../users/users.module";
import { ARTokensService } from "./ar-tokens.service";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";
import { EmailChangeService } from "./email-change.service";
import { EmailVerificationService } from "./email-verification.service";
import { PasswordResetService } from "./password-reset.service";
import { SessionSerializer } from "./serializers/session.serializer";
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
    UsersModule,
    forwardRef(() => ApiTokensModule),
    JwtModule.registerAsync({
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        secret: configService.get<string>("JWT_API_TOKEN_SECRET"),
      }),
    }),
    OrganizationsModule,
    PrismaModule,
    EmailModule,
    PassportModule.register({ session: false }),
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
    ARTokensService,
    // Serializers
    SessionSerializer,
  ],
})
export class AuthModule {}
