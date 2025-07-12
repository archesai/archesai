import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { SuccessDocumentSchema } from '#http/schemas/success-document.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createIndividualResponseSchema = (
  resourceObjectSchema: TObject,
  entityKey: string
) => {
  return Type.Composite(
    [
      Type.Object({
        data: resourceObjectSchema
      }),
      SuccessDocumentSchema
    ],
    {
      // $id: `${toTitleCaseNoSpaces(entityKey)}IndividualResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} Individual response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Individual Response`
    }
  )
}
