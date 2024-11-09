import { applyDecorators } from "@nestjs/common";
import {
  ApiBadRequestResponse,
  ApiOkResponse,
  ApiOperation,
  ApiResponse,
} from "@nestjs/swagger";

import { Roles } from "../../auth/decorators/roles.decorator";
import { ApiPaginatedResponse } from "./paginated.decorator";

export enum Operation {
  CREATE = "CREATE",
  DELETE = "DELETE",
  FIND_ALL = "FIND_ALL",
  GET = "GET",
  UPDATE = "UPDATE",
}

export function ApiCrudOperation<TEntity>(
  operation: Operation,
  entityName: string,
  entityType: new (...args: any[]) => TEntity,
  isAdmin: boolean,
  aggregates?: new (...args: any[]) => any
) {
  let summary = `${
    operation.charAt(0) + operation.slice(1).toLowerCase()
  } a ${entityName}`;
  const defaultResponses = [
    ApiResponse({ description: "Unauthorized", status: 401 }),
    ApiResponse({ description: "Not Found", status: 404 }),
  ];

  const specificResponses = [];

  if (isAdmin) {
    specificResponses.push(
      ApiResponse({ description: "Forbidden", status: 403 })
    );
  }

  switch (operation) {
    case "CREATE":
      summary = `Create a new ${entityName}`;
      specificResponses.push(
        ApiResponse({
          description: `Successfully created a new ${entityName}`,
          status: 201,
          type: entityType,
        })
      );
      specificResponses.push(
        ApiBadRequestResponse({
          description: `Bad request when creating the ${entityName}`,
        })
      );
      break;
    case "DELETE":
      specificResponses.push(
        ApiOkResponse({
          description: `Successfully deleted the ${entityName}`,
        })
      );
      break;
    case "FIND_ALL":
      summary = `Get all ${entityName}s`;
      specificResponses.push(ApiPaginatedResponse(entityType, aggregates));
      break;
    case "GET":
      specificResponses.push(
        ApiOkResponse({
          description: `Successfully retrieved the ${entityName}`,
          type: entityType,
        })
      );
      break;
    case "UPDATE":
      specificResponses.push(
        ApiOkResponse({
          description: `Successfully updated the ${entityName}`,
          type: entityType,
        })
      );
      break;
  }

  const description =
    summary +
    (isAdmin ? ". ADMIN ONLY." : ". USER and ADMIN can access this endpoint.");

  return applyDecorators(
    Roles(...(isAdmin ? ["ADMIN"] : ["ADMIN", "USER"])),
    ApiOperation({
      description,
      summary,
    }),
    ...defaultResponses,
    ...specificResponses
  );
}
