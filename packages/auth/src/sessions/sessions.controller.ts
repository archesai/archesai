import { Type } from '@sinclair/typebox'

import type {
  ArchesApiRequest,
  ArchesApiResponse,
  Controller,
  HttpInstance
} from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef } from '@archesai/domain'

import type { SessionsService } from '#sessions/sessions.service'

import { CreateAccountRequestSchema } from '#accounts/dto/create-account.req.dto'
import { AuthenticatedGuard } from '#auth/guards/authenticated.guard'
import { LocalAuthGuard } from '#auth/guards/local-auth.guard'

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
    return { ...request.user }
  }

  public async login(
    _request: ArchesApiRequest,
    _reply: ArchesApiResponse
  ): Promise<void> {
    // The LocalAuthGuard will handle the login
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
            204: LegacyRef(ArchesApiNoContentResponseSchema),
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
            200: Type.Object({}),
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
