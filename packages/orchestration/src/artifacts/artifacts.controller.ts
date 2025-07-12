import type { Controller } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  ARTIFACT_ENTITY_KEY,
  ArtifactEntitySchema,
  CreateArtifactDtoSchema,
  UpdateArtifactDtoSchema
} from '@archesai/schemas'

import type { ArtifactsService } from '#artifacts/artifacts.service'

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
      CreateArtifactDtoSchema,
      UpdateArtifactDtoSchema,
      artifactsService
    )
  }
}
