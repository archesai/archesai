import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type { InvitationEntity } from '@archesai/schemas'

import {
  ArchesApiForbiddenResponseSchema,
  ArchesApiNotFoundResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  AuthenticatedGuard,
  BaseController,
  toTitleCase
} from '@archesai/core'
import {
  CreateInvitationDtoSchema,
  INVITATION_ENTITY_KEY,
  InvitationEntitySchema,
  LegacyRef,
  UpdateInvitationDtoSchema
} from '@archesai/schemas'

import type { InvitationsService } from '#invitations/invitations.service'

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
      CreateInvitationDtoSchema,
      UpdateInvitationDtoSchema,
      invitationsService
    )
    this.invitationsService = invitationsService
  }

  public async accept(request: ArchesApiRequest) {
    return this.toIndividualResponse(
      await this.invitationsService.accept('FIXME', request.user!)
    )
  }

  public override registerRoutes(app: HttpInstance) {
    super.registerRoutes(app)
    app.post(
      `/${this.entityKey}/:id/accept`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: 'Accept an invitation',
          operationId: 'acceptInvitation',
          // params: Type.Object({
          //   id: Type.String()
          // }),
          response: {
            200: this.invididualResponseSchema,
            401: LegacyRef(ArchesApiUnauthorizedResponseSchema),
            403: LegacyRef(ArchesApiForbiddenResponseSchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: 'Accept an invitation',
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.accept.bind(this)
    )
  }
}
