import type { Logger } from '@archesai/core'
import type { ContentEntity } from '@archesai/domain'

import type { ContentService } from '#content/content.service'

// A tool run processor should take in the runId, the input contents, a logger, and the content service
// It should return the output contents
export type Transformer = (
  runId: string,
  inputs: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  ...args: unknown[]
) => Promise<ContentEntity[]>
