import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type { UpdateEmailVerificationDto } from '@archesai/schemas'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  AuthenticatedGuard,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef, UpdateEmailVerificationDtoSchema } from '@archesai/schemas'

import type { EmailVerificationService } from '#email-verification/email-verification.service'

/**
 * Controller for managing email verifications.
 */
export class EmailVerificationController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly emailVerificationService: EmailVerificationService

  constructor(emailVerificationService: EmailVerificationService) {
    this.emailVerificationService = emailVerificationService
  }

  public async confirm(
    request: ArchesApiRequest & { body: UpdateEmailVerificationDto }
  ): Promise<void> {
    await this.emailVerificationService.confirm(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/email-verification/confirm`,
      {
        schema: {
          body: UpdateEmailVerificationDtoSchema,
          description: 'This endpoint will confirm your e-mail with a token',
          operationId: 'confirmEmailVerification',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema),
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: 'Confirm e-mail verification',
          tags: ['Email Verification']
        }
      },
      this.confirm.bind(this)
    )

    app.post(
      `/auth/email-verification/request`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description:
            'This endpoint will send an e-mail verification link to you. ADMIN ONLY.',
          operationId: 'requestEmailVerification',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          security: [{ bearerAuth: [] }], // ✅ add this line
          summary: 'Request e-mail verification',
          tags: ['Email Verification']
        }
      },
      this.request.bind(this)
    )
  }

  public async request(request: ArchesApiRequest): Promise<void> {
    return this.emailVerificationService.request({
      email: request.user!.email,
      userId: request.user!.id
    })
  }
}
