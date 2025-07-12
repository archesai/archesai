import type { Controller } from '@archesai/core'
import type { LabelEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateLabelDtoSchema,
  LABEL_ENTITY_KEY,
  LabelEntitySchema,
  UpdateLabelDtoSchema
} from '@archesai/schemas'

import type { LabelsService } from '#labels/labels.service'

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
      CreateLabelDtoSchema,
      UpdateLabelDtoSchema,
      labelsService
    )
  }
}
