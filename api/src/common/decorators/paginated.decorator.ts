import { applyDecorators, Type } from '@nestjs/common'
import { ApiOkResponse, getSchemaPath } from '@nestjs/swagger'

import { AggregateFieldResult, Metadata } from '../dto/paginated.dto'

export const ApiPaginatedResponse = <TModel extends Type<any>>(model: TModel) => {
  return applyDecorators(
    ApiOkResponse({
      description: 'Successfully returned paginated results',
      schema: {
        allOf: [
          {
            properties: {
              aggregates: {
                items: { $ref: getSchemaPath(AggregateFieldResult) },
                type: 'array'
              },
              metadata: {
                $ref: getSchemaPath(Metadata)
              },
              results: {
                items: { $ref: getSchemaPath(model) },
                type: 'array'
              }
            }
          }
        ],
        title: `PaginatedResponseOf${model.name}`
      }
    })
  )
}
