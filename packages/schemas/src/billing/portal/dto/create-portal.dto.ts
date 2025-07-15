import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreatePortalDtoSchema: TObject<{
  organizationId: TString
}> = Type.Object({
  organizationId: Type.String()
})

export type CreatePortalDto = Static<typeof CreatePortalDtoSchema>
