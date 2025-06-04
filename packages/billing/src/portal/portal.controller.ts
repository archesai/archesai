import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import { IS_CONTROLLER } from '@archesai/core'

import type { CreatePortalRequest } from '#portal/dto/create-portal.req.dto'
import type { PortalService } from '#portal/portal.service'

import { CreatePortalRequestSchema } from '#portal/dto/create-portal.req.dto'
import { PortalResourceSchema } from '#portal/dto/portal.res.dto'

/**
 * Controller for billing portal.
 */
export class PortalController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly portalService: PortalService

  constructor(portalService: PortalService) {
    this.portalService = portalService
  }

  public async create(
    request: ArchesApiRequest & { body: CreatePortalRequest }
  ) {
    return this.portalService.create(request.body)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/billing/portal`,
      {
        schema: {
          body: CreatePortalRequestSchema,
          description: 'Create a new portal',
          operationId: 'createPortal',
          response: {
            201: {
              description: 'The created portal',
              schema: PortalResourceSchema
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
