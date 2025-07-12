import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreatePortalDtoSchema = Type.Object({
  organizationId: Type.String()
})

export type CreatePortalDto = Static<typeof CreatePortalDtoSchema>
