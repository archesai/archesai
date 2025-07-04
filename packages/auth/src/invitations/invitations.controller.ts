import { Type } from '@sinclair/typebox'

import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type { InvitationEntity } from '@archesai/domain'

import {
  ArchesApiForbiddenResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  BaseController,
  toTitleCase
} from '@archesai/core'
import {
  INVITATION_ENTITY_KEY,
  InvitationEntitySchema,
  LegacyRef
} from '@archesai/domain'

import type { InvitationsService } from '#invitations/invitations.service'

import { CreateInvitationRequestSchema } from '#invitations/dto/create-invitation.req.dto'
import { UpdateInvitationRequestSchema } from '#invitations/dto/update-invitation.req.dto'

/**
 * Controller for handling invitations.
 */
export class InvitationsController
  extends BaseController<InvitationEntity>
  implements Controller
{
  private readonly invitationsService: InvitationsService
  constructor(invitationsService: InvitationsService) {
    super(
      INVITATION_ENTITY_KEY,
      InvitationEntitySchema,
      CreateInvitationRequestSchema,
      UpdateInvitationRequestSchema,
      invitationsService
    )
    this.invitationsService = invitationsService
  }

  public async accept(request: ArchesApiRequest & { params: { id: string } }) {
    return this.toIndividualResponse(
      await this.invitationsService.accept(request.params.id, request.user!)
    )
  }

  public override registerRoutes(app: HttpInstance) {
    super.registerRoutes(app)
    app.post(
      `/${INVITATION_ENTITY_KEY}/:id/accept`,
      {
        schema: {
          description: 'Accept an invitation',
          operationId: 'acceptInvitation',
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: this.IndividualEntityResponseSchema,
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema),
            403: LegacyRef(ArchesApiForbiddenResponseSchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: 'Accept an invitation',
          tags: [toTitleCase(INVITATION_ENTITY_KEY)]
        }
      },
      this.accept.bind(this)
    )
  }
}
