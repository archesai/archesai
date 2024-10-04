import { HttpModule } from "@nestjs/axios";
import { forwardRef, Module } from "@nestjs/common";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { JwtModule } from "@nestjs/jwt";

import { ApiTokensModule } from "../api-tokens/api-tokens.module";
import { FirebaseModule } from "../firebase/firebase.module";
import { OrganizationsModule } from "../organizations/organizations.module";
import { UsersModule } from "../users/users.module";
import { AuthController } from "./auth.controller";
import { AuthService } from "./auth.service";
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
    FirebaseModule,
    forwardRef(() => ApiTokensModule),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        secret: configService.get<string>("JWT_API_TOKEN_SECRET"),
      }),
    }),
    OrganizationsModule,
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
  ],
})
export class AuthModule {}
