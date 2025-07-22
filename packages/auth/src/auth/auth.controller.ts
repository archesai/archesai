import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import { UnauthorizedException } from '@archesai/core'
import {
  CreateAccountDtoSchema,
  DocumentSchemaFactory,
  NoContentResponseSchema,
  UnauthorizedResponseSchema,
  UserEntitySchema
} from '@archesai/schemas'

import type { AuthService } from '#auth/auth.service'

import { getHeaders, setHeaders } from '#utils/headers'

export interface AuthPluginOptions {
  authService: AuthService
  // organizationsService: OrganizationsService
}

export const authPlugin: FastifyPluginAsyncZod<AuthPluginOptions> = async (
  app,
  { authService }
) => {
  app.post(
    `/sign-up`,
    {
      schema: {
        body: CreateAccountDtoSchema,
        description: `This endpoint will register you with your e-mail and password`,
        operationId: 'register',
        response: {
          201: DocumentSchemaFactory(UserEntitySchema),
          401: UnauthorizedResponseSchema
        },
        summary: `Register`,
        tags: ['Authentication']
      }
    },
    async (req, res) => {
      const {
        headers,
        response: { user }
      } = await authService.signUpEmail({
        body: {
          email: req.body.email,
          name: req.body.name,
          password: req.body.password
        },
        returnHeaders: true
      })

      // create an organization for the user
      // const organization = await organizationsService.create({
      // body: {
      //   name: `${user.email}'s Organization`,
      //   slug: user.email,
      //   userId: user.id
      // },
      // })

      // set the active organization for the user
      // fiix me

      setHeaders(headers, res)
      return {
        data: {
          user
        }
      }
    }
  )

  app.post(
    `/sign-in`,
    {
      schema: {
        body: CreateAccountDtoSchema.pick({
          email: true,
          password: true
        }),
        description: `This endpoint will log you in with your e-mail and password`,
        operationId: 'login',
        response: {
          200: DocumentSchemaFactory(UserEntitySchema),
          401: UnauthorizedResponseSchema
        },
        summary: `Login`,
        tags: ['Authentication']
      }
    },
    async (req, res) => {
      const { headers, response } = await authService.signInEmail({
        body: {
          email: req.body.email,
          password: req.body.password
        },
        returnHeaders: true
      })
      setHeaders(headers, res)
      return {
        data: response.user
      }
    }
  )

  app.post(
    `/sign-out`,
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
    async (req, res) => {
      const reqHeaders = getHeaders(req.headers)
      const {
        headers: resHeaders,
        response: { success }
      } = await authService.signOut({
        headers: reqHeaders,
        returnHeaders: true
      })
      if (!success) {
        throw new UnauthorizedException(
          'You are not logged in or your session has expired.'
        )
      }
      setHeaders(resHeaders, res)
      return null
    }
  )

  await Promise.resolve()
}
