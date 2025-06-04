import { Type } from '@sinclair/typebox'

import { LabelEntitySchema } from '@archesai/domain'

export const CreateLabelRequestSchema = Type.Object({
  name: LabelEntitySchema.properties.name
})
