import { Test, TestingModule } from '@nestjs/testing'
import { ConfirmationTokenWithNewPasswordDto } from '@/src/auth/dto/confirmation-token-with-new-password.dto'
import { EmailRequestDto } from '@/src/auth/dto/email-request.dto'
import { CookiesDto } from '@/src/auth/dto/token.dto'
import { AuthService } from '@/src/auth/services/auth.service'
import { PasswordResetService } from '@/src/auth/services/password-reset.service'
import { PasswordResetController } from '../controllers/password-reset.controller'
import { Response } from 'express'

describe('PasswordResetController', () => {
  let passwordResetController: PasswordResetController
  let passwordResetService: PasswordResetService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [PasswordResetController],
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
          provide: PasswordResetService,
          useValue: {
            confirm: jest.fn(),
            request: jest.fn()
          }
        }
      ]
    }).compile()

    passwordResetController = module.get<PasswordResetController>(
      PasswordResetController
    )
    passwordResetService =
      module.get<PasswordResetService>(PasswordResetService)
  })

  it('should be defined', () => {
    expect(passwordResetController).toBeDefined()
  })

  describe('passwordResetConfirm', () => {
    it('should confirm password reset', async () => {
      const res: Response = {} as Response
      const dto: ConfirmationTokenWithNewPasswordDto = {
        newPassword: 'new-password',
        token: 'test-token'
      }
      const result = new CookiesDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(passwordResetService, 'confirm').mockResolvedValue(result)

      expect(
        await passwordResetController.passwordResetConfirm(dto, res)
      ).toEqual(result)
    })
  })

  describe('passwordResetRequest', () => {
    it('should request password reset', async () => {
      const dto: EmailRequestDto = { email: 'test@example.com' }
      jest.spyOn(passwordResetService, 'request').mockResolvedValue(undefined)

      expect(
        await passwordResetController.passwordResetRequest(dto)
      ).toBeUndefined()
    })
  })
})
