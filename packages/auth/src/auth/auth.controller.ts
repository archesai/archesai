import type { FastifyReply, FastifyRequest } from 'fastify'
import type { FastifyPluginCallbackZod } from 'fastify-type-provider-zod'

import type { SessionEntity, UserEntity } from '@archesai/schemas'

import {
  BetterAuthSessionSchema,
  CreateAccountDtoSchema,
  CreateEmailChangeDtoSchema,
  CreatePasswordResetDtoSchema,
  NoContentResponseSchema,
  NotFoundResponseSchema,
  // SessionEntitySchema,
  UnauthorizedResponseSchema,
  UpdateEmailChangeDtoSchema,
  UpdateEmailVerificationDtoSchema,
  UpdatePasswordResetDtoSchema
  // UserEntitySchema
} from '@archesai/schemas'

import type { AuthService } from '#auth/auth.service'

export interface AuthPluginOptions {
  authService: AuthService
}

declare module 'fastify' {
  interface FastifyInstance {
    authHandler: (
      req: FastifyRequest,
      reply: FastifyReply,
      beforeSend?: (
        response: Response,
        responseText: null | string
      ) => Promise<void>
    ) => Promise<void>
  }
}

export const authPlugin: FastifyPluginCallbackZod<AuthPluginOptions> = (
  app,
  { authService },
  done
) => {
  // Optional: Add helper methods to fastify instance
  app.decorate(
    'authHandler',
    async (
      req: FastifyRequest,
      reply: FastifyReply,
      beforeSend?: (
        response: Response,
        responseText: null | string
      ) => Promise<void>
    ) => {
      try {
        // Reusable auth handler logic that can be called from other routes
        const url = new URL(
          req.url,
          `http://${req.headers.host?.toString() ?? ''}`
        )

        const headers = new Headers()
        Object.entries(req.headers).forEach(([key, value]) => {
          if (value) headers.append(key, value.toString())
        })

        const formattedRequest = new Request(url.toString(), {
          body: req.body ? JSON.stringify(req.body) : undefined,
          headers,
          method: req.method
        })

        const response = await authService.handler(formattedRequest)

        // Get response text once
        const responseText = response.body ? await response.text() : null

        // Forward response to client
        reply.status(response.status)
        response.headers.forEach((value, key) => {
          reply.header(key, value)
        })

        // Run callback if provided
        if (beforeSend) {
          await beforeSend(response, responseText)
        }

        reply.send(responseText)
      } catch (err) {
        app.log.error('Authentication Error:', err)
        reply.status(500).send({
          errocode: 'AUTH_FAILURE',
          error: 'Internal authentication error'
        })
      }
    }
  )

  app.post(
    `/api/auth/sign-in/email`,
    {
      schema: {
        body: CreateAccountDtoSchema,
        description: `This endpoint will log you in with your e-mail and password`,
        operationId: 'login',
        response: {
          200: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema
        },
        summary: `Login`,
        tags: ['Authentication']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/sign-out`,
    {
      schema: {
        description: `This endpoint will log you out of the current session`,
        operationId: 'logout',
        response: {
          204: NoContentResponseSchema,
          401: UnauthorizedResponseSchema
        },
        summary: `Logout`,
        tags: ['Authentication']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.get(
    `/api/auth/get-session`,
    {
      schema: {
        description: `This endpoint will return the current session information`,
        operationId: 'getSession',
        response: {
          200: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema
        },
        summary: `Get Session`,
        tags: ['Authentication']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/sign-up/email`,
    {
      schema: {
        body: CreateAccountDtoSchema,
        description: `This endpoint will register you with your e-mail and password`,
        operationId: 'register',
        response: {
          201: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema
        },
        summary: `Register`,
        tags: ['Registration']
      }
    },
    (req, res) => {
      return app.authHandler(req, res, async (response, responseText) => {
        // Create organization after successful signup
        if (response.status === 200 && responseText) {
          try {
            const userData = JSON.parse(responseText) as {
              session: SessionEntity
              user: UserEntity
            }
            if (userData.user.id) {
              const organization = await authService.createOrganization({
                body: {
                  name: `${userData.user.email}'s Organization`,
                  slug: userData.user.email,
                  userId: userData.user.id
                }
              })
              if (!organization) {
                throw new Error('Failed to create organization')
              }
              await authService.setActiveOrganization({
                body: {
                  organizationId: organization.id
                }
              })
            }
          } catch (orgError) {
            console.error(
              'Failed to create organization after signup:',
              orgError
            )
          }
        }
      })
    }
  )

  app.post(
    `/api/auth/verify-email`,
    {
      schema: {
        body: UpdateEmailVerificationDtoSchema,
        description: 'This endpoint will confirm your e-mail with a token',
        operationId: 'confirmEmailVerification',
        response: {
          200: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Confirm e-mail verification',
        tags: ['Email Verification']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/send-verification-email`,
    {
      schema: {
        description:
          'This endpoint will send an e-mail verification link to you. ADMIN ONLY.',
        operationId: 'requestEmailVerification',
        response: {
          204: NoContentResponseSchema
        },
        security: [{ bearerAuth: [] }], // ✅ add this line
        summary: 'Request e-mail verification',
        tags: ['Email Verification']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/forgot-password`,
    {
      schema: {
        body: UpdatePasswordResetDtoSchema,
        description:
          'This endpoint will confirm your password change with a token',
        operationId: 'confirmPasswordReset',
        response: {
          204: NoContentResponseSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Confirm password reset',
        tags: ['Password Reset']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/reset-password`,
    {
      schema: {
        body: CreatePasswordResetDtoSchema,
        description: 'This endpoint will request a password reset link',
        operationId: 'requestPasswordReset',
        response: {
          204: NoContentResponseSchema
        },
        summary: 'Request password reset',
        tags: ['Password Reset']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/change-email`,
    {
      schema: {
        body: UpdateEmailChangeDtoSchema,
        description:
          'This endpoint will confirm your e-mail change with a token',
        operationId: 'confirmEmailChange',
        response: {
          204: NoContentResponseSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Confirm e-mail change',
        tags: ['Email Change']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  app.post(
    `/api/auth/email-change/request`,
    {
      schema: {
        body: CreateEmailChangeDtoSchema,
        description:
          'This endpoint will request your e-mail change with a token',
        operationId: 'requestEmailChange',
        response: {
          204: NoContentResponseSchema
        },
        summary: 'Request e-mail change',
        tags: ['Email Change']
      }
    },
    (req, res) => {
      return app.authHandler(req, res)
    }
  )

  done()
}
