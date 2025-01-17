import { HttpModule } from '@nestjs/axios'
import { forwardRef, Module } from '@nestjs/common'
import { JwtModule } from '@nestjs/jwt'
import { PassportModule } from '@nestjs/passport'

import { EmailModule } from '@/src/email/email.module'
import { OrganizationsModule } from '@/src/organizations/organizations.module'
import { PrismaModule } from '@/src/prisma/prisma.module'
import { UsersModule } from '@/src/users/users.module'
import { SessionSerializer } from '@/src/auth/serializers/session.serializer'
import { ARTokensService } from '@/src/auth/services/ar-tokens.service'
import { AuthService } from '@/src/auth/services/auth.service'
import { EmailChangeService } from '@/src/auth/services/email-change.service'
import { EmailVerificationService } from '@/src/auth/services/email-verification.service'
import { PasswordResetService } from '@/src/auth/services/password-reset.service'
import { ApiKeyStrategy } from '@/src/auth/strategies/api-key.strategy'
import { FirebaseStrategy } from '@/src/auth/strategies/firebase.strategy'
import { JwtStrategy } from '@/src/auth/strategies/jwt.strategy'
import { LocalStrategy } from '@/src/auth/strategies/local-strategy'
import { TwitterStrategy } from '@/src/auth/strategies/twitter.strategy'
import { ConfigService } from '@/src/config/config.service'
import { EmailChangeController } from '@/src/auth/controllers/email-change.controller'
import { EmailVerificationController } from '@/src/auth/controllers/email-verification.controller'
import { ProvidersController } from '@/src/auth/controllers/providers.controller'
import { PasswordResetController } from '@/src/auth/controllers/password-reset.controller'
import { AuthController } from '@/src/auth/controllers/auth.controller'
import { ApiTokensModule } from '@/src/api-tokens/api-tokens.module'
import { ConfigModule } from '@/src/config/config.module'

@Module({
  controllers: [
    AuthController,
    EmailChangeController,
    EmailVerificationController,
    ProvidersController,
    PasswordResetController
  ],
  exports: [AuthService],
  imports: [
    HttpModule,
    UsersModule,
    forwardRef(() => ApiTokensModule),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        secret: configService.get('jwt.secret')
      })
    }),
    OrganizationsModule,
    PrismaModule,
    EmailModule,
    PassportModule.register({ session: false })
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
    SessionSerializer
  ]
})
export class AuthModule {}
