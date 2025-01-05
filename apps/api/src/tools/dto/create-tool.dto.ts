import { PickType } from '@nestjs/swagger'

import { ToolEntity } from '../entities/tool.entity'

export class CreateToolDto extends PickType(ToolEntity, [
  'name',
  'description',
  'inputType',
  'outputType',
  'toolBase'
] as const) {}
