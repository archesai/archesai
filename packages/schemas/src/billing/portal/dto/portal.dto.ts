import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const PortalDtoSchema: TObject<{
  url: TString
}> = Type.Object({
  url: Type.String({
    description: 'The URL that will bring you to the necessary Stripe page'
  })
})

export type PortalDto = Static<typeof PortalDtoSchema>
