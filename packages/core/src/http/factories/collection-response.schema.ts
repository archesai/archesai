import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { SuccessDocumentSchema } from '#http/schemas/success-document.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createCollectionResponseSchema = <T extends TObject>(
  resourceObjectSchema: T,
  entityKey: string
) => {
  return Type.Composite(
    [
      Type.Object({
        // data: Type.Array(
        //   Type.Unsafe<StaticDecode<T>>(Type.Ref(resourceObjectSchema.$id!))
        // )
        data: Type.Array(resourceObjectSchema)
      }),
      SuccessDocumentSchema
    ],
    {
      // $id: `${toTitleCaseNoSpaces(entityKey)}CollectionResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} collection response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Collection Response`
    }
  )
}
