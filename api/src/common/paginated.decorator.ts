import { applyDecorators, Type } from "@nestjs/common";
import { ApiOkResponse, getSchemaPath } from "@nestjs/swagger";

import { _PaginatedDto } from "./paginated.dto";

export const ApiPaginatedResponse = <
  TModel extends Type<any>,
  TAggregateModel extends Type<any> = any,
>(
  model: TModel,
  aggregateModel?: TAggregateModel
) => {
  return applyDecorators(
    ApiOkResponse({
      description: "Successfully returned paginated results",
      schema: {
        allOf: [
          { $ref: getSchemaPath(_PaginatedDto) },
          {
            properties: {
              results: {
                items: { $ref: getSchemaPath(model) },
                type: "array",
              },
              ...(aggregateModel
                ? {
                    aggregates: {
                      $ref: getSchemaPath(aggregateModel),
                    },
                  }
                : {}),
            },
          },
        ],
        title: `PaginatedResponseOf${model.name}`,
      },
    })
  );
};
