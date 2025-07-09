import type { Controller } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import { ARTIFACT_ENTITY_KEY, ArtifactEntitySchema } from '@archesai/schemas'

import type { ArtifactsService } from '#artifacts/artifacts.service'

import { CreateArtifactRequestSchema } from '#artifacts/dto/create-artifact.req.dto'
import { UpdateArtifactRequestSchema } from '#artifacts/dto/update-artifact.req.dto'

/**
 * Controller for content.
 */
export class ArtifactsController
  extends BaseController<ArtifactEntity>
  implements Controller
{
  constructor(artifactsService: ArtifactsService) {
    super(
      ARTIFACT_ENTITY_KEY,
      ArtifactEntitySchema,
      CreateArtifactRequestSchema,
      UpdateArtifactRequestSchema,
      artifactsService
    )
  }
}
