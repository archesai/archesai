import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type {
  CreatePasswordResetDto,
  UpdatePasswordResetDto
} from '@archesai/schemas'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  IS_CONTROLLER
} from '@archesai/core'
import {
  CreatePasswordResetDtoSchema,
  LegacyRef,
  UpdatePasswordResetDtoSchema
} from '@archesai/schemas'

import type { PasswordResetService } from '#password-reset/password-reset.service'

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
    request: ArchesApiRequest & { body: UpdatePasswordResetDto }
  ): Promise<void> {
    return this.passwordResetService.confirm(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/password-reset/confirm`,
      {
        schema: {
          body: UpdatePasswordResetDtoSchema,
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
          body: CreatePasswordResetDtoSchema,
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
    request: ArchesApiRequest & {
      body: CreatePasswordResetDto
    }
  ): Promise<void> {
    return this.passwordResetService.request(request.body)
  }
}
