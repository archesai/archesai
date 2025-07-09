import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef } from '@archesai/schemas'

import type { CreatePasswordResetRequest } from '#password-reset/dto/create-password-reset.req.dto'
import type { UpdatePasswordResetRequest } from '#password-reset/dto/update-password-reset.req.dto'
import type { PasswordResetService } from '#password-reset/password-reset.service'

import { CreatePasswordResetRequestSchema } from '#password-reset/dto/create-password-reset.req.dto'
import { UpdatePasswordResetRequestSchema } from '#password-reset/dto/update-password-reset.req.dto'

/**
 * Controller for password reset.
 */
export class PasswordResetController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly passwordResetService: PasswordResetService

  constructor(passwordResetService: PasswordResetService) {
    this.passwordResetService = passwordResetService
  }

  public async confirm(
    request: ArchesApiRequest & { body: UpdatePasswordResetRequest }
  ): Promise<void> {
    return this.passwordResetService.confirm(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/password-reset/confirm`,
      {
        schema: {
          body: UpdatePasswordResetRequestSchema,
          description:
            'This endpoint will confirm your password change with a token',
          operationId: 'confirmPasswordReset',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema),
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: 'Confirm password reset',
          tags: ['Password Reset']
        }
      },
      this.confirm.bind(this)
    )

    app.post(
      `/auth/password-reset/request`,
      {
        schema: {
          body: CreatePasswordResetRequestSchema,
          description: 'This endpoint will request a password reset link',
          operationId: 'requestPasswordReset',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          summary: 'Request password reset',
          tags: ['Password Reset']
        }
      },
      this.request.bind(this)
    )
  }

  public async request(
    request: ArchesApiRequest & { body: CreatePasswordResetRequest }
  ): Promise<void> {
    return this.passwordResetService.request(request.body)
  }
}
