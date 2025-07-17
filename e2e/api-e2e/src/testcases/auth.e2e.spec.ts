import request from 'supertest'

import type { HttpInstance } from '@archesai/core'

describe('Auth Module E2E Tests', () => {
  let app: HttpInstance
  let accessToken: string
  let emailChangeToken = ''
  let emailVerification: null | string = null
  let passwordResetToken: null | string = null

  const userCredentials = {
    email: 'auth-e2e-test@archesai.com',
    password: 'Password123!'
  }

  const newEmail = 'new-email@archesai.com'

  beforeAll(async () => {
    app = await createApp()

    // Mock EmailService
    const emailService = app.get(EmailService)
    jest
      .spyOn(emailService, 'sendMail')
      .mockImplementation(async ({ html }) => {
        if (!html || typeof html !== 'string') {
          return Promise.resolve()
        }
        const tokenMatch = /token=([a-zA-Z0-9]+)/.exec(html)
        if (tokenMatch) {
          const token = tokenMatch[1]!
          // Determine token type based on URL path or another identifier
          if (html.toString().includes('email-change')) {
            emailChangeToken = token
          } else if (html.toString().includes('email-verification')) {
            emailVerification = token
          } else if (html.toString().includes('password-reset')) {
            passwordResetToken = token
          }
        }
        return Promise.resolve()
      })

    await app.init()

    // Register user
    const tokenDto = await registerUser(app, userCredentials)
    accessToken = tokenDto.accessToken
  })

  afterAll(async () => {
    await app.close()
  })

  describe('Authentication', () => {
    it('Should protect private endpoints like /user', async () => {
      // Attempt to access protected route without token
      const resUnauthorized = await request(app.getHttpServer()).get('/user')
      expect(resUnauthorized.status).toBe(401)

      // Access protected route with valid token
      const resAuthorized = await request(app.getHttpServer())
        .get('/user')
        .set('Authorization', `Bearer ${accessToken}`)
      expect(resAuthorized.status).toBe(200)
      expect(resAuthorized.body.email).toBe(userCredentials.email)
      expect(resAuthorized).toSatisfyApiSpec()
    })

    it('Should deactivate user', async () => {
      await deactivateUser(app, accessToken)

      // Verify that the user is deactivated by attempting to access /user
      const resAfterDeactivation = await request(app.getHttpServer())
        .get('/user')
        .set('Authorization', `Bearer ${accessToken}`)
      expect(resAfterDeactivation.status).toBe(403)

      // Reactivate the user for email change tests
    })
  })

  describe('Email Change', () => {
    beforeAll(async () => {
      const userRepository = app.get(UserRepository)
      await userRepository.update('id', {
        deactivated: false
      })
      // where: { email: userCredentials.email }
    })

    it('Should request an email change', async () => {
      const emailRequestDto: EmailRequestDto = {
        email: newEmail
      }
      const res = await request(app.getHttpServer())
        .post('/auth/email-change/request')
        .set('Authorization', `Bearer ${accessToken}`)
        .send(emailRequestDto)
      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
      expect(emailChangeToken).not.toBeNull()
    })

    it('Should throw an error if email change token is invalid', async () => {
      const invalidTokenDto: VerificationDto = {
        token: 'invalid-token'
      }
      const res = await request(app.getHttpServer())
        .post('/auth/email-change/confirm')
        .send(invalidTokenDto)
      expect(res.status).toBe(400)
    })

    it('Should confirm the email change', async () => {
      expect(emailChangeToken).not.toBeNull()
      const validTokenDto: VerificationDto = {
        token: emailChangeToken
      }
      const res = await request(app.getHttpServer())
        .post('/auth/email-change/confirm')
        .send(validTokenDto)

      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
    })

    it('Should update the email', async () => {
      const user = await getUser(app, accessToken)
      expect(user.email).toBe(newEmail)
    })
  })

  describe('Email Verification', () => {
    it('Should default users to not email verified', async () => {
      const user = await getUser(app, accessToken)
      expect(user.emailVerified).toBe(false)
    })

    it('Should send an email verification link', async () => {
      const res = await request(app.getHttpServer())
        .post('/auth/email-verification/request')
        .set('Authorization', `Bearer ${accessToken}`)

      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
      expect(emailVerification).not.toBeNull()
    })

    it('Should confirm email verification', async () => {
      expect(emailVerification).not.toBeNull()
      const validTokenDto: VerificationDto = {
        token: emailVerification ?? ''
      }
      const res = await request(app.getHttpServer())
        .post('/auth/email-verification/confirm')
        .send(validTokenDto)

      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
      expect(res.body.accessToken).toBeDefined()

      const user = await getUser(app, accessToken)
      expect(user.emailVerified).toBe(true)
    })

    it('Should throw 400 if email verification token is invalid', async () => {
      const invalidTokenDto: VerificationDto = {
        token: 'invalid-token'
      }
      const res = await request(app.getHttpServer())
        .post('/auth/email-verification/confirm')
        .send(invalidTokenDto)

      expect(res.status).toBe(400)
    })
  })

  describe('Password Reset', () => {
    it('Should request a password reset', async () => {
      const passwordResetRequestDto: EmailRequestDto = {
        email: newEmail // Use the updated email
      }
      const res = await request(app.getHttpServer())
        .post('/auth/password-reset/request')
        .send(passwordResetRequestDto)

      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
      expect(passwordResetToken).not.toBeNull()
    })

    it('Should not indicate whether the email exists', async () => {
      const passwordResetRequestDto: EmailRequestDto = {
        email: 'non-existent@archesai.com'
      }
      const res = await request(app.getHttpServer())
        .post('/auth/password-reset/request')
        .send(passwordResetRequestDto)

      // Even if the email does not exist, respond with 201 to prevent email enumeration
      expect(res.status).toBe(201)
      expect(res).toSatisfyApiSpec()
    })

    it('Should throw an error if password reset token is invalid', async () => {
      const invalidTokenDto: VerificationDto & { newPassword: string } = {
        newPassword: 'NewPassword123!',
        token: 'invalid-token'
      }
      const res = await request(app.getHttpServer())
        .post('/auth/password-reset/confirm')
        .send(invalidTokenDto)

      expect(res.status).toBe(400)
    })

    it('Should confirm the password reset', async () => {
      expect(passwordResetToken).not.toBeNull()
      const validResetDto: VerificationDto & { newPassword: string } = {
        newPassword: 'NewPassword123!',
        token: passwordResetToken ?? ''
      }
      const res = await request(app.getHttpServer())
        .post('/auth/password-reset/confirm')
        .send(validResetDto)

      expect(res.status).toBe(201)
      expect(res.body.accessToken).toBeDefined()

      // Attempt to login with the new password
      const loginRes = await request(app.getHttpServer())
        .post('/auth/login')
        .send({
          email: newEmail,
          password: 'NewPassword123!'
        })
      expect(loginRes.status).toBe(201)
      expect(loginRes.body.accessToken).toBeDefined()
      expect(loginRes).toSatisfyApiSpec()
    })

    it('Should prevent token reuse after successful password reset', async () => {
      expect(passwordResetToken).not.toBeNull()
      const validResetDto: VerificationDto & { newPassword: string } = {
        newPassword: 'AnotherNewPassword123!',
        token: passwordResetToken ?? ''
      }
      const res1 = await request(app.getHttpServer())
        .post('/auth/password-reset/confirm')
        .send(validResetDto)

      expect(res1.status).toBe(400)
    })

    it('Should handle token expiration correctly', async () => {
      // Request another password reset to get a fresh token
      const passwordResetRequestDto: EmailRequestDto = {
        email: newEmail
      }
      await request(app.getHttpServer())
        .post('/auth/password-reset/request')
        .send(passwordResetRequestDto)
        .expect(201)

      expect(passwordResetToken).not.toBeNull()

      // Manually expire the token
      const VerificationRepository = app.get(VerificationRepository)
      await VerificationRepository.update('id', {
        expires: new Date(Date.now() - 1000) // Set to past
      })

      const expiredResetDto: VerificationDto & { newPassword: string } = {
        newPassword: 'ExpiredPassword123!',
        token: passwordResetToken ?? ''
      }
      const res = await request(app.getHttpServer())
        .post('/auth/password-reset/confirm')
        .send(expiredResetDto)

      expect(res.status).toBe(400)

      // Clean up: remove the expired token
      await VerificationRepository.delete('id')

      passwordResetToken = null
    })
  })
})
