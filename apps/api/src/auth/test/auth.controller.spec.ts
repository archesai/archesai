import { createRandomUser } from '@/prisma/factories/user.factory'
import { Test, TestingModule } from '@nestjs/testing'
import { Response } from 'express'
import { Request } from 'express'

import { AuthController } from '@/src/auth/controllers/auth.controller'

import { LoginDto } from '@/src/auth/dto/login.dto'
import { RegisterDto } from '@/src/auth/dto/register.dto'
import { CookiesDto } from '@/src/auth/dto/token.dto'
import { AuthService } from '@/src/auth/services/auth.service'
import { EmailChangeService } from '@/src/auth/services/email-change.service'
import { EmailVerificationService } from '@/src/auth/services/email-verification.service'
import { PasswordResetService } from '@/src/auth/services/password-reset.service'

describe('AuthController', () => {
  let authController: AuthController
  let authService: AuthService

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
  })

  it('should be defined', () => {
    expect(authController).toBeDefined()
  })

  describe('login', () => {
    it('should login user', async () => {
      const dto: LoginDto = { email: 'test@example.com', password: 'password' }
      const user = createRandomUser()
      const result = new CookiesDto({
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

  describe('refreshToken', () => {
    it('should refresh token', async () => {
      const req: Partial<Request> = {
        signedCookies: { 'archesai.refreshToken': 'test-refresh-token' }
      }
      const res: Response = {} as Response
      const result = new CookiesDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(authService, 'refreshAccessToken').mockResolvedValue(result)
      jest.spyOn(authService, 'setCookies').mockResolvedValue(undefined)

      expect(await authController.refreshToken(req as Request, res)).toEqual(
        result
      )
    })
  })

  describe('register', () => {
    it('should register user', async () => {
      const dto: RegisterDto = {
        email: 'test@example.com',
        password: 'password'
      }
      const res: Response = {} as Response
      const user = createRandomUser()
      const result = new CookiesDto({
        accessToken: 'accessToken',
        refreshToken: 'refreshToken'
      })
      jest.spyOn(authService, 'register').mockResolvedValue(user)
      jest.spyOn(authService, 'login').mockResolvedValue(result)

      expect(await authController.register(dto, res)).toEqual(result)
    })
  })
})
