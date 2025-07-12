import type {
  ArchesApiRequest,
  ArchesApiResponse,
  Controller,
  HttpInstance
} from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  AuthenticatedGuard,
  IS_CONTROLLER,
  UnauthorizedException
} from '@archesai/core'
import {
  CreateAccountDtoSchema,
  LegacyRef,
  UserEntitySchema
} from '@archesai/schemas'

/**
 * Controller for managing authentication.
 */
export class SessionsController implements Controller {
  public readonly [IS_CONTROLLER] = true

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/login`,
      {
        preValidation: [LocalAuthGuard(app)],
        schema: {
          body: CreateAccountDtoSchema,
          description: `This endpoint will log you in with your e-mail and password`,
          operationId: 'login',
          response: {
            201: UserEntitySchema,
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema)
          },
          summary: `Login`,
          tags: ['Sessions']
        }
      },
      (req) => {
        return req.user!
      }
    )

    app.post(
      `/auth/logout`,
      {
        schema: {
          description: `This endpoint will log you out of the current session`,
          operationId: 'logout',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema),
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema)
          },
          summary: `Logout`,
          tags: ['Sessions']
        }
      },
      (req) => {
        return req.logOut()
      }
    )

    app.get(
      `/auth/session`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `This endpoint will return the current session information`,
          operationId: 'getSession',
          response: {
            200: UserEntitySchema,
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema)
          },
          summary: `Get Session`,
          tags: ['Sessions']
        }
      },
      (req) => {
        return req.user!
      }
    )
  }
}

export function LocalAuthGuard(app: HttpInstance) {
  return async function (
    req: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<void> {
    await new Promise<void>((resolve, reject) => {
      const handler = req.passport.authenticate(
        ['local'],
        { session: true },
        async (authReq, _authRes, err, user) => {
          if (err) {
            reject(err)
            return
          }
          if (!user) {
            reject(new Error('Unauthorized'))
            return
          }

          try {
            await authReq.logIn(user)
            resolve()
          } catch {
            reject(new UnauthorizedException())
          }
        }
      )

      handler.call(app, req, reply)
    })
  }
}
