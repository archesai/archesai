import type {
  ArchesApiRequest,
  ArchesApiResponse,
  Controller,
  HttpInstance
} from '@archesai/core'
import type { UserEntity } from '@archesai/domain'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  AuthenticatedGuard,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef, UserEntitySchema } from '@archesai/domain'

import type { SessionsService } from '#sessions/sessions.service'

import { CreateAccountRequestSchema } from '#accounts/dto/create-account.req.dto'

/**
 * Controller for managing authentication.
 */
export class SessionsController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly sessionsService: SessionsService

  constructor(sessionsService: SessionsService) {
    this.sessionsService = sessionsService
  }

  public getSession(request: ArchesApiRequest) {
    return { ...request.user } as UserEntity
  }

  public login(
    request: ArchesApiRequest,
    _reply: ArchesApiResponse
  ): UserEntity {
    // The LocalAuthGuard will handle the login
    return request.user!
  }

  public async logout(
    request: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<void> {
    await this.sessionsService.logout(request, reply)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/login`,
      {
        preValidation: [LocalAuthGuard(app)],
        schema: {
          body: CreateAccountRequestSchema,
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
      this.login.bind(this)
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
      this.logout.bind(this)
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
      this.getSession.bind(this)
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
          } catch (err) {
            reject(err as Error)
          }
        }
      )

      handler.call(app, req, reply)
    })
  }
}
