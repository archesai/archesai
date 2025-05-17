import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const createEntitySchema = (
  entityKey: string,
  attributesSchema: TObject,
  relationshipsSchema?: TObject
) => {
  return Type.Object(
    {
      attributes: attributesSchema,
      id: Type.String({
        description: `The ${entityKey} id`
      }),
      ...(relationshipsSchema
        ? {
            relationships: relationshipsSchema
          }
        : {}),
      type: Type.Literal(entityKey)
    },
    {
      $id: `${entityKey}Entity`,
      description: `The ${entityKey} entity`,
      title: `${entityKey} Entity`
    }
  )
}
