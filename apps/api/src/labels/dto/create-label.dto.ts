import { PickType } from '@nestjs/swagger'

import { LabelEntity } from '../entities/label.entity'

export class CreateLabelDto extends PickType(LabelEntity, ['name'] as const) {}
