import { createRandomUser } from '@/prisma/factories/user.factory'
import { Test, TestingModule } from '@nestjs/testing'
import { EmailChangeController } from '@/src/auth/controllers/email-change.controller'
import { ConfirmationTokenDto } from '@/src/auth/dto/confirmation-token.dto'
import { EmailRequestDto } from '@/src/auth/dto/email-request.dto'
import { CookiesDto } from '@/src/auth/dto/token.dto'
import { AuthService } from '@/src/auth/services/auth.service'
import { EmailChangeService } from '@/src/auth/services/email-change.service'

describe('EmailChangeController', () => {
  let emailChangeController: EmailChangeController
  let emailChangeService: EmailChangeService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [EmailChangeController],
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
          provide: EmailChangeService,
          useValue: {
            confirm: jest.fn(),
            request: jest.fn()
          }
        }
      ]
    }).compile()

    emailChangeController = module.get<EmailChangeController>(
      EmailChangeController
    )
    emailChangeService = module.get<EmailChangeService>(EmailChangeService)
  })

  it('should be defined', () => {
    expect(emailChangeController).toBeDefined()
  })

  describe('emailChangeConfirm', () => {
    it('should confirm email change', async () => {
      const dto: ConfirmationTokenDto = { token: 'confirmationToken' }
      const result = new CookiesDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(emailChangeService, 'confirm').mockResolvedValue(result)
      expect(await emailChangeController.emailChangeConfirm(dto)).toEqual(
        result
      )
    })
  })

  describe('emailChangeRequest', () => {
    it('should request email change', async () => {
      const user = createRandomUser()
      const dto: EmailRequestDto = { email: 'test@example.com' }
      jest.spyOn(emailChangeService, 'request').mockResolvedValue(undefined)
      expect(
        await emailChangeController.emailChangeRequest(user, dto)
      ).toBeUndefined()
    })
  })
})
