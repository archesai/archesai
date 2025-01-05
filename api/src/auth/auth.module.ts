import { HttpModule } from '@nestjs/axios'
import { forwardRef, Module } from '@nestjs/common'
import { JwtModule } from '@nestjs/jwt'
import { PassportModule } from '@nestjs/passport'

import { ApiTokensModule } from '../api-tokens/api-tokens.module'
import { EmailModule } from '../email/email.module'
import { OrganizationsModule } from '../organizations/organizations.module'
import { PrismaModule } from '../prisma/prisma.module'
import { UsersModule } from '../users/users.module'
import { AuthController } from './auth.controller'
import { SessionSerializer } from './serializers/session.serializer'
import { ARTokensService } from './services/ar-tokens.service'
import { AuthService } from './services/auth.service'
import { EmailChangeService } from './services/email-change.service'
import { EmailVerificationService } from './services/email-verification.service'
import { PasswordResetService } from './services/password-reset.service'
import { ApiKeyStrategy } from './strategies/api-key.strategy'
import { FirebaseStrategy } from './strategies/firebase.strategy'
import { JwtStrategy } from './strategies/jwt.strategy'
import { LocalStrategy } from './strategies/local-strategy'
import { TwitterStrategy } from './strategies/twitter.strategy'
import { ArchesConfigService } from '../config/config.service'

@Module({
  controllers: [AuthController],
  exports: [AuthService],
  imports: [
    HttpModule,
    UsersModule,
    forwardRef(() => ApiTokensModule),
    JwtModule.registerAsync({
      inject: [ArchesConfigService],
      useFactory: async (configService: ArchesConfigService) => ({
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
