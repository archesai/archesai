import { randomBytes } from 'node:crypto'

import type { ConfigService, EmailService } from '@archesai/core'
import type { VerificationTokenType } from '@archesai/schemas'

import { UnauthorizedException } from '@archesai/core'

import type { HashingService } from '#hashing/hashing.service'
import type { VerificationTokenRepository } from '#verification-tokens/verification-token.repository'

/**
 * Service for verification tokens.
 */
export class VerificationTokensService {
  private readonly configService: ConfigService
  private readonly emailService: EmailService
  private readonly hashingService: HashingService
  private readonly verificationTokenRepository: VerificationTokenRepository

  constructor(
    configService: ConfigService,
    emailService: EmailService,
    hashingService: HashingService,
    verificationTokenRepository: VerificationTokenRepository
  ) {
    this.configService = configService
    this.emailService = emailService
    this.hashingService = hashingService
    this.verificationTokenRepository = verificationTokenRepository
  }

  public async create(
    _type: VerificationTokenType,
    userId: string,
    expiresInHours: number,
    overrides?: { newEmail?: string }
  ): Promise<string> {
    await this.verificationTokenRepository.deleteMany({
      filter: {}
    })

    const token = randomBytes(32).toString('hex')
    const expiresAt = new Date()
    expiresAt.setHours(expiresAt.getHours() + expiresInHours)
    await this.verificationTokenRepository.create({
      expires: expiresAt.toISOString(),
      identifier: userId,
      newEmail: overrides?.newEmail ?? '',
      token: await this.hashingService.hashPassword(token)
    })

    return token
  }

  public async sendNotification(
    type: VerificationTokenType,
    token: string,
    email: string,
    subject: string,
    htmlFunction: (link: string, email: string) => string
  ): Promise<void> {
    const lowercaseType = type.toString().replace('_', '-').toLowerCase()
    const changeEmailLink = `${this.configService.get('platform.host')}/confirm?type=${lowercaseType}&token=${token}`
    const html = htmlFunction(changeEmailLink, email)
    await this.emailService.sendMail({
      html,
      subject,
      to: email
    })
  }

  public async verify(_type: VerificationTokenType, token: string) {
    const verificationToken = await this.verificationTokenRepository.findFirst({
      filter: {
        token: {
          equals: token
        }
      }
    })

    if (
      await this.hashingService.verifyPassword(token, verificationToken.token)
    ) {
      if (new Date(verificationToken.expires) < new Date()) {
        throw new UnauthorizedException('Token has expired.')
      }
      await this.verificationTokenRepository.delete(verificationToken.id)
      return {
        newEmail: verificationToken.newEmail,
        userId: verificationToken.identifier
      }
    } else {
      throw new UnauthorizedException('Invalid or expired token.')
    }
  }
}
