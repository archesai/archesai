import { createRandomUser } from '@/prisma/factories/user.factory'
import { Test, TestingModule } from '@nestjs/testing'
import { Response } from 'express'
import { Request } from 'express'

import { AuthController } from '../auth.controller'
import { ConfirmationTokenWithNewPasswordDto } from '../dto/confirmation-token-with-new-password.dto'
import { ConfirmationTokenDto } from '../dto/confirmation-token.dto'
import { EmailRequestDto } from '../dto/email-request.dto'
import { LoginDto } from '../dto/login.dto'
import { RegisterDto } from '../dto/register.dto'
import { TokenDto } from '../dto/token.dto'
import { AuthService } from '../services/auth.service'
import { EmailChangeService } from '../services/email-change.service'
import { EmailVerificationService } from '../services/email-verification.service'
import { PasswordResetService } from '../services/password-reset.service'

describe('AuthController', () => {
  let authController: AuthController
  let authService: AuthService
  let emailChangeService: EmailChangeService
  let emailVerificationService: EmailVerificationService
  let passwordResetService: PasswordResetService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [AuthController],
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
        },
        {
          provide: EmailVerificationService,
          useValue: {
            confirm: jest.fn(),
            request: jest.fn()
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

    authController = module.get<AuthController>(AuthController)
    authService = module.get<AuthService>(AuthService)
    emailChangeService = module.get<EmailChangeService>(EmailChangeService)
    emailVerificationService = module.get<EmailVerificationService>(
      EmailVerificationService
    )
    passwordResetService =
      module.get<PasswordResetService>(PasswordResetService)
  })

  it('should be defined', () => {
    expect(authController).toBeDefined()
  })

  describe('emailChangeConfirm', () => {
    it('should confirm email change', async () => {
      const dto: ConfirmationTokenDto = { token: 'confirmationToken' }
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(emailChangeService, 'confirm').mockResolvedValue(result)

      expect(await authController.emailChangeConfirm(dto)).toEqual(result)
    })
  })

  describe('emailChangeRequest', () => {
    it('should request email change', async () => {
      const user = createRandomUser()
      const dto: EmailRequestDto = { email: 'test@example.com' }
      jest.spyOn(emailChangeService, 'request').mockResolvedValue(undefined)

      expect(await authController.emailChangeRequest(user, dto)).toBeUndefined()
    })
  })

  describe('emailVerificationConfirm', () => {
    it('should confirm email verification', async () => {
      const dto: ConfirmationTokenDto = { token: 'confirmationToken' }
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(emailVerificationService, 'confirm').mockResolvedValue(result)

      expect(await authController.emailVerificationConfirm(dto)).toEqual(result)
    })
  })

  describe('emailVerificationRequest', () => {
    it('should request email verification', async () => {
      const user = createRandomUser()
      jest
        .spyOn(emailVerificationService, 'request')
        .mockResolvedValue(undefined)
      expect(
        await authController.emailVerificationRequest(user)
      ).toBeUndefined()
    })
  })

  describe('login', () => {
    it('should login user', async () => {
      const dto: LoginDto = { email: 'test@example.com', password: 'password' }
      const user = createRandomUser()
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      const res: Response = {} as Response
      jest.spyOn(authService, 'login').mockResolvedValue(result)
      jest.spyOn(authService, 'setCookies').mockResolvedValue(undefined)

      expect(await authController.login(dto, user, res)).toEqual(result)
    })
  })

  describe('logout', () => {
    it('should logout user', async () => {
      const res: Response = {} as Response
      jest.spyOn(authService, 'removeCookies').mockResolvedValue(undefined)

      expect(await authController.logout(res)).toBeUndefined()
    })
  })

  describe('passwordResetConfirm', () => {
    it('should confirm password reset', async () => {
      const dto: ConfirmationTokenWithNewPasswordDto = {
        newPassword: 'new-password',
        token: 'test-token'
      }
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(passwordResetService, 'confirm').mockResolvedValue(result)

      expect(await authController.passwordResetConfirm(dto)).toEqual(result)
    })
  })

  describe('passwordResetRequest', () => {
    it('should request password reset', async () => {
      const dto: EmailRequestDto = { email: 'test@example.com' }
      jest.spyOn(passwordResetService, 'request').mockResolvedValue(undefined)

      expect(await authController.passwordResetRequest(dto)).toBeUndefined()
    })
  })

  describe('refreshToken', () => {
    it('should refresh token', async () => {
      const refreshToken = 'test-refresh-token'
      const req: Partial<Request> = {
        signedCookies: { 'archesai.refreshToken': 'test-refresh-token' }
      }
      const res: Response = {} as Response
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(authService, 'refreshAccessToken').mockResolvedValue(result)
      jest.spyOn(authService, 'setCookies').mockResolvedValue(undefined)

      expect(
        await authController.refreshToken(refreshToken, req as Request, res)
      ).toEqual(result)
    })
  })

  describe('register', () => {
    it('should register user', async () => {
      const dto: RegisterDto = {
        email: 'test@example.com',
        password: 'password'
      }
      const user = createRandomUser()
      const result = new TokenDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(authService, 'register').mockResolvedValue(user)
      jest.spyOn(authService, 'login').mockResolvedValue(result)

      expect(await authController.register(dto)).toEqual(result)
    })
  })
})
