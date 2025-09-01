import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'
import type { z } from 'zod'

import type { BaseEntity, SearchQuery } from '@archesai/schemas'

// import { AuthenticatedGuard } from '#http/guards/authenticated.guard'
import {
  createSearchQuerySchema,
  DocumentCollectionSchemaFactory,
  DocumentSchemaFactory,
  IdParamsSchema,
  NotFoundResponseSchema
} from '@archesai/schemas'

import type { BaseService } from '#common/base-service'

import { capitalize } from '#utils/capitalize'
import { singularize } from '#utils/pluralize'
import { toCamelCase, toTitleCase, vf } from '#utils/strings'

export interface CrudPluginOptions<
  TEntity extends BaseEntity = BaseEntity,
  TInsert = unknown,
  TSelect extends TEntity = TEntity,
  TCreateSchema extends z.ZodType<TInsert> = z.ZodType<TInsert>,
  TUpdateSchema extends z.ZodType<Partial<TInsert>> = z.ZodType<
    Partial<TInsert>
  >
> {
  createSchema?: TCreateSchema
  enableBulkOperations?: boolean
  entityKey: string
  entitySchema: z.ZodObject
  prefix: string
  service: BaseService<TEntity, TInsert, TSelect>
  tags?: string[]
  updateSchema?: TUpdateSchema
}

export const crudPlugin: FastifyPluginAsyncZod<CrudPluginOptions> = async (
  app,
  {
    createSchema,
    entityKey,
    entitySchema,
    service,
    tags = [entityKey],
    updateSchema
  }
) => {
  const baseRouteOptions = {
    // preValidation: [AuthenticatedGuard()],
    schema: {
      security: [{ bearerAuth: [] }],
      tags
    }
  }

  const searchQuerySchema = createSearchQuerySchema(entitySchema, entityKey)

  if (createSchema) {
    // POST /entity - Create single entity
    app.post(
      '',
      {
        ...baseRouteOptions,
        schema: {
          ...baseRouteOptions.schema,
          body: createSchema,
          description: `Create a new ${singularize(entityKey)}`,
          operationId:
            'create' + capitalize(toCamelCase(singularize(entityKey))),
          response: {
            201: DocumentSchemaFactory(entitySchema)
            // 400: { $ref: 'error-document' },
            // 401: { $ref: 'unauthorized-response' },
            // 403: { $ref: 'forbidden-response' }
          },
          summary: `Create a new ${singularize(entityKey)}`,
          tags: [toTitleCase(entityKey)]
        }
      },
      async (req) => {
        return {
          data: await service.create(req.body)
        }
      }
    )
  }

  // GET /entity - Find many entities
  app.get(
    '',
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Find many ${entityKey}`,
        operationId: 'findMany' + capitalize(toCamelCase(entityKey)),
        querystring: searchQuerySchema,
        response: {
          200: DocumentCollectionSchemaFactory(entitySchema)
        },
        summary: `Find many ${entityKey}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    async (request) => {
      const results = await service.findMany(
        request.query as SearchQuery<unknown>
      )
      return {
        data: results.data,
        meta: {
          total: results.count
        }
      }
    }
  )

  // DELETE /entity/:id - Delete single entity
  app.delete(
    `/:id`,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Delete a${vf(entityKey)} ${singularize(entityKey)}`,
        operationId: 'delete' + capitalize(toCamelCase(singularize(entityKey))),
        params: IdParamsSchema,
        response: {
          200: DocumentSchemaFactory(entitySchema),
          404: NotFoundResponseSchema
        },
        summary: `Delete a${vf(entityKey)} ${singularize(entityKey)}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    async (request) => {
      return {
        data: await service.delete(request.params.id)
      }
    }
  )

  // GET /entity/:id - Find single entity by ID
  app.get(
    `/:id`,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Find a${vf(entityKey)} ${singularize(entityKey)}`,
        operationId: 'getOne' + capitalize(toCamelCase(singularize(entityKey))),
        params: IdParamsSchema,
        response: {
          200: DocumentSchemaFactory(entitySchema),
          404: NotFoundResponseSchema
        },
        summary: `Find a${vf(entityKey)} ${singularize(entityKey)}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    async (request) => {
      return {
        data: await service.findOne(request.params.id)
      }
    }
  )

  // PATCH /entity/:id - Update single entity
  if (updateSchema) {
    app.patch(
      `/:id`,
      {
        ...baseRouteOptions,
        schema: {
          ...baseRouteOptions.schema,
          body: updateSchema,
          description: `Update a${vf(entityKey)} ${singularize(entityKey)}`,
          operationId:
            'update' + capitalize(toCamelCase(singularize(entityKey))),
          params: IdParamsSchema,
          response: {
            200: DocumentSchemaFactory(entitySchema),
            404: NotFoundResponseSchema
          },
          summary: `Update a${vf(entityKey)} ${singularize(entityKey)}`,
          tags: [toTitleCase(entityKey)]
        }
      },
      async (request) => {
        return {
          data: await service.update(request.params.id, request.body)
        }
      }
    )
  }

  await Promise.resolve()
}

// // POST /entity/bulk - Create multiple entities (optional)
// if (enableBulkOperations) {
//   app.post(
//     `${prefix}/bulk`,
//     {
//       ...baseRouteOptions,
//       schema: {
//         ...baseRouteOptions.schema,
//         body: {
//           items: createSchema,
//           type: 'array'
//         },
//         response: {
//           200: {
//             properties: {
//               count: { type: 'number' },
//               data: { items: entitySchema, type: 'array' }
//             },
//             type: 'object'
//           },
//           400: { $ref: 'error-document' },
//           401: { $ref: 'unauthorized-response' },
//           403: { $ref: 'forbidden-response' }
//         },
//         summary: `Create multiple ${entityKey}s`
//       }
//     },
//     async (request) => {
//       return service.createMany(request.body)
//     }
//   )
// }

//   if (enableBulkOperations) {
//   app.delete(
//     `${prefix}/bulk`,
//     {
//       ...baseRouteOptions,
//       schema: {
//         ...baseRouteOptions.schema,
//         body: {
//           properties: {
//             query: {
//               properties: {
//                 limit: { type: 'number' },
//                 offset: { type: 'number' },
//                 search: { type: 'string' }
//               },
//               type: 'object'
//             }
//           },
//           required: ['query'],
//           type: 'object'
//         },
//         response: {
//           200: {
//             properties: {
//               count: { type: 'number' },
//               data: { items: entitySchema, type: 'array' }
//             },
//             type: 'object'
//           },
//           400: { $ref: 'error-document' },
//           401: { $ref: 'unauthorized-response' },
//           403: { $ref: 'forbidden-response' }
//         },
//         summary: `Delete multiple ${entityKey}s`
//       }
//     },
//     async (req) => {
//       return service.deleteMany(req.query)
//     }
//   )
// }

// if (enableBulkOperations) {
//   app.patch(
//     `${prefix}/bulk`,
//     {
//       ...baseRouteOptions,
//       schema: {
//         ...baseRouteOptions.schema,
//         body: {
//           properties: {
//             data: updateSchema,
//             query: {
//               properties: {
//                 limit: { type: 'number' },
//                 offset: { type: 'number' },
//                 search: { type: 'string' }
//               },
//               type: 'object'
//             }
//           },
//           required: ['data'],
//           type: 'object'
//         },
//         response: {
//           200: {
//             properties: {
//               count: { type: 'number' },
//               data: { items: entitySchema, type: 'array' }
//             },
//             type: 'object'
//           },
//           400: { $ref: 'error-document' },
//           401: { $ref: 'unauthorized-response' },
//           403: { $ref: 'forbidden-response' }
//         },
//         summary: `Update multiple ${entityKey}s`
//       }
//     },
//     async (request) => {
//       const { data, query } = request.body as {
//         data: any
//         query?: SearchQuery<any>
//       }
//       return service.updateMany(data, query || {})
//     }
//   )
// }
