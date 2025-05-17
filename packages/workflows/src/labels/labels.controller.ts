import type { Controller } from '@archesai/core'
import type { LabelEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { LABEL_ENTITY_KEY, LabelEntitySchema } from '@archesai/domain'

import type { LabelsService } from '#labels/labels.service'

import { CreateLabelRequestSchema } from '#labels/dto/create-label.req.dto'
import { UpdateLabelRequestSchema } from '#labels/dto/update-label.req.dto'

/**
 * Controller for labels.
 */
export class LabelsController
  extends BaseController<LabelEntity>
  implements Controller
{
  constructor(labelsService: LabelsService) {
    super(
      LABEL_ENTITY_KEY,
      LabelEntitySchema,
      CreateLabelRequestSchema,
      UpdateLabelRequestSchema,
      labelsService
    )
  }
}
