import type { Logger } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/domain'

import type { ArtifactsService } from '#artifacts/artifacts.service'

// A tool run processor should take in the runId, the input contents, a logger, and the content service
// It should return the output contents
export type Transformer = (
  runId: string,
  inputs: ArtifactEntity[],
  logger: Logger,
  artifactsService: ArtifactsService,
  ...args: unknown[]
) => Promise<ArtifactEntity[]>
