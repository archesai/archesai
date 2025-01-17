import { createRandomUser } from '@/prisma/factories/user.factory'
import { Test, TestingModule } from '@nestjs/testing'

import { EmailVerificationController } from '@/src/auth/controllers/email-verification.controller'
import { ConfirmationTokenDto } from '@/src/auth/dto/confirmation-token.dto'

import { CookiesDto } from '@/src/auth/dto/token.dto'
import { AuthService } from '@/src/auth/services/auth.service'
import { EmailVerificationService } from '@/src/auth/services/email-verification.service'
import { Response } from 'express'

describe('EmailVerificationController', () => {
  let emailVerificationController: EmailVerificationController
  let emailVerificationService: EmailVerificationService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [EmailVerificationController],
      providers: [
        {
          provide: AuthService,
          useValue: {
            login: jest.fn(),
            refreshAccessToken: jest.fn(),
            register: jest.fn(),
            removeCookies: jest.fn(),
            setCookies: jest.fn()
          }
        },

        {
          provide: EmailVerificationService,
          useValue: {
            confirm: jest.fn(),
            request: jest.fn()
          }
        }
      ]
    }).compile()

    emailVerificationController = module.get<EmailVerificationController>(
      EmailVerificationController
    )
    emailVerificationService = module.get<EmailVerificationService>(
      EmailVerificationService
    )
  })

  it('should be defined', () => {
    expect(emailVerificationController).toBeDefined()
  })

  describe('emailVerificationConfirm', () => {
    it('should confirm email verification', async () => {
      const res: Response = {} as Response
      const dto: ConfirmationTokenDto = { token: 'confirmationToken' }
      const result = new CookiesDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(emailVerificationService, 'confirm').mockResolvedValue(result)
      expect(
        await emailVerificationController.emailVerificationConfirm(dto, res)
      ).toEqual(result)
    })
  })

  describe('emailVerificationRequest', () => {
    it('should request email verification', async () => {
      const user = createRandomUser()
      jest
        .spyOn(emailVerificationService, 'request')
        .mockResolvedValue(undefined)
      expect(
        await emailVerificationController.emailVerificationRequest(user)
      ).toBeUndefined()
    })
  })
})
