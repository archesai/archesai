import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  AuthenticatedGuard,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef } from '@archesai/schemas'

import type { CreateEmailChangeRequest } from '#email-change/dto/create-email-change-request.dto'
import type { UpdateEmailChangeRequest } from '#email-change/dto/update-email-change-request.dto'
import type { EmailChangeService } from '#email-change/email-change.service'

import { CreateEmailChangeRequestSchema } from '#email-change/dto/create-email-change-request.dto'
import { UpdateEmailChangeRequestSchema } from '#email-change/dto/update-email-change-request.dto'

/**
 * Controller for managing email changes.
 */
export class EmailChangeController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly emailChangeService: EmailChangeService

  constructor(emailChangeService: EmailChangeService) {
    this.emailChangeService = emailChangeService
  }

  public async confirm(
    request: ArchesApiRequest & { body: UpdateEmailChangeRequest }
  ): Promise<void> {
    await this.emailChangeService.confirm(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/email-change/confirm`,
      {
        schema: {
          body: UpdateEmailChangeRequestSchema,
          description:
            'This endpoint will confirm your e-mail change with a token',
          operationId: 'confirmEmailChange',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema),
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: 'Confirm e-mail change',
          tags: ['Email Change']
        }
      },
      this.confirm.bind(this)
    )

    app.post(
      `/auth/email-change/request`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          body: CreateEmailChangeRequestSchema,
          description:
            'This endpoint will request your e-mail change with a token',
          operationId: 'requestEmailChange',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          summary: 'Request e-mail change',
          tags: ['Email Change']
        }
      },
      this.request.bind(this)
    )
  }

  public async request(
    request: ArchesApiRequest & {
      body: CreateEmailChangeRequest
    }
  ): Promise<void> {
    return this.emailChangeService.request(request.body)
  }
}
