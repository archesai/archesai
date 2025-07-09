import { Type } from '@sinclair/typebox'

import { LabelEntitySchema } from '@archesai/schemas'

export const CreateLabelRequestSchema = Type.Object({
  name: LabelEntitySchema.properties.name
})
