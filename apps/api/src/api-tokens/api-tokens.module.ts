import { Module } from '@nestjs/common'
import { JwtModule } from '@nestjs/jwt'

import { PrismaModule } from '../prisma/prisma.module'
import { ApiTokenRepository } from './api-token.repository'
import { ApiTokensController } from './api-tokens.controller'
import { ApiTokensService } from './api-tokens.service'
import { AuthModule } from '../auth/auth.module'

@Module({
  controllers: [ApiTokensController],
  exports: [ApiTokensService],
  imports: [PrismaModule, JwtModule, AuthModule],
  providers: [ApiTokensService, ApiTokenRepository]
})
export class ApiTokensModule {}
