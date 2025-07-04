import type { StaticDecode } from '@sinclair/typebox'

import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  IS_CONTROLLER
} from '@archesai/core'
import { LegacyRef } from '@archesai/domain'

import type { RegistrationService } from '#registration/registration.service'

import { CreateAccountRequestSchema } from '#accounts/dto/create-account.req.dto'

/**
 * Controller for managing registration.
 */
export class RegistrationController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly registrationService: RegistrationService

  constructor(registrationService: RegistrationService) {
    this.registrationService = registrationService
  }

  public async register(
    request: ArchesApiRequest & {
      body: StaticDecode<typeof CreateAccountRequestSchema>
    }
  ): Promise<void> {
    const user = await this.registrationService.register(
      request.body.email,
      request.body.password
    )
    await request.logIn(user)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/register`,
      {
        schema: {
          body: CreateAccountRequestSchema,
          description: `This endpoint will register you with your e-mail and password`,
          operationId: 'register',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema),
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema)
          },
          summary: `Register`,
          tags: ['Registration']
        }
      },
      this.register.bind(this)
    )
  }
}
