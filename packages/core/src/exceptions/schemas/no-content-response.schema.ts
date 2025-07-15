import type { TNull } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const NoContentResponseSchema: TNull = Type.Null({
  $id: 'NoContentResponse',
  description: 'No Content',
  title: 'No Content'
})
