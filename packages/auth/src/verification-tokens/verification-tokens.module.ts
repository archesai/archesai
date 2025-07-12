import type { ModuleMetadata } from '@archesai/core'
import type { VerificationTokenEntity } from '@archesai/schemas'

import {
  ConfigModule,
  ConfigService,
  createModule,
  DatabaseModule,
  DatabaseService,
  EmailModule,
  EmailService
} from '@archesai/core'

import { HashingModule } from '#hashing/hashing.module'
import { HashingService } from '#hashing/hashing.service'
import { VerificationTokenRepository } from '#verification-tokens/verification-token.repository'
import { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

export const VerificationTokensModuleDefinition: ModuleMetadata = {
  exports: [VerificationTokensService],
  imports: [ConfigModule, DatabaseModule, EmailModule, HashingModule],
  providers: [
    {
      inject: [
        ConfigService,
        EmailService,
        HashingService,
        VerificationTokenRepository
      ],
      provide: VerificationTokensService,
      useFactory: (
        configService: ConfigService,
        emailService: EmailService,
        hashingService: HashingService,
        verificationTokenRepository: VerificationTokenRepository
      ) =>
        new VerificationTokensService(
          configService,
          emailService,
          hashingService,
          verificationTokenRepository
        )
    },
    {
      inject: [DatabaseService],
      provide: VerificationTokenRepository,
      useFactory: (databaseService: DatabaseService<VerificationTokenEntity>) =>
        new VerificationTokenRepository(databaseService)
    }
  ]
}

export const VerificationTokensModule = (() =>
  createModule(
    class VerificationTokensModule {},
    VerificationTokensModuleDefinition
  ))()
