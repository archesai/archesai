import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ApiResponseSchema } from '#http/schemas/api-response.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createIndividualResponseSchema = (
  ResourceObjectSchema: TObject,
  entityKey: string
) => {
  return Type.Composite(
    [
      Type.Object({
        data: ResourceObjectSchema
      }),
      ApiResponseSchema
    ],
    {
      // $id: `${toTitleCaseNoSpaces(entityKey)}IndividualResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} Individual response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Individual Response`
    }
  )
}
