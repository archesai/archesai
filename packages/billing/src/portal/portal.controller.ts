import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type { CreatePortalDto } from '@archesai/schemas'

import { IS_CONTROLLER } from '@archesai/core'
import { CreatePortalDtoSchema, PortalDtoSchema } from '@archesai/schemas'

import type { PortalService } from '#portal/portal.service'

/**
 * Controller for billing portal.
 */
export class PortalController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly portalService: PortalService

  constructor(portalService: PortalService) {
    this.portalService = portalService
  }

  public async create(request: ArchesApiRequest & { body: CreatePortalDto }) {
    return this.portalService.create(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/billing/portal`,
      {
        schema: {
          body: CreatePortalDtoSchema,
          description: 'Create a new portal',
          operationId: 'createPortal',
          response: {
            201: {
              description: 'The created portal',
              schema: PortalDtoSchema
            }
          },
          summary: 'Create a new portal',
          tags: ['Billing']
        }
      },
      this.create.bind(this)
    )
  }
}
