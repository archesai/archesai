import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreatePortalRequestSchema = Type.Object({
  organizationId: Type.String()
})

export type CreatePortalRequest = Static<typeof CreatePortalRequestSchema>
