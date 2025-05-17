import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ApiResponseSchema } from '#http/schemas/api-response.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createCollectionResponseSchema = <T extends TObject>(
  ResourceObjectSchema: T,
  entityKey: string
) => {
  return Type.Composite(
    [
      Type.Object({
        // data: Type.Array(
        //   Type.Unsafe<StaticDecode<T>>(Type.Ref(ResourceObjectSchema.$id!))
        // )
        data: Type.Array(ResourceObjectSchema)
      }),
      ApiResponseSchema
    ],
    {
      // $id: `${toTitleCaseNoSpaces(entityKey)}CollectionResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} collection response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Collection Response`
    }
  )
}
