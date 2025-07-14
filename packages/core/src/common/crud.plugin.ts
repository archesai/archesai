import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { BaseEntity, TObject, TSchema } from '@archesai/schemas'

import { LegacyRef, Type } from '@archesai/schemas'

import type { BaseService } from '#common/base-service'

import { NotFoundResponseSchema } from '#exceptions/schemas/not-found-response.schema'
import { createSearchQuerySchema } from '#http/dto/search-query.dto'
// import { AuthenticatedGuard } from '#http/guards/authenticated.guard'
import { DocumentSchemaFactory } from '#http/schemas/document.schema'
import { capitalize } from '#utils/capitalize'
import { singularize } from '#utils/pluralize'
import { toCamelCase, toTitleCase, vf } from '#utils/strings'

export interface CrudPluginOptions<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity,
  TCreateSchema extends TSchema = TSchema,
  TUpdateSchema extends TSchema = TSchema
> {
  createSchema: TCreateSchema
  enableBulkOperations?: boolean
  entityKey: string
  entitySchema: TSchema
  prefix: string
  service: BaseService<TEntity, TInsert, TSelect>
  tags?: string[]
  updateSchema: TUpdateSchema
}

export const crudPlugin: FastifyPluginAsyncTypebox<
  CrudPluginOptions<BaseEntity, unknown, BaseEntity>
> = async (
  app,
  {
    createSchema,
    // enableBulkOperations = false,
    entityKey,
    entitySchema,
    prefix,
    service,
    tags = [entityKey],
    updateSchema
  }
  // eslint-disable-next-line @typescript-eslint/require-await
) => {
  const baseRouteOptions = {
    // preValidation: [AuthenticatedGuard()],
    schema: {
      security: [{ bearerAuth: [] }],
      tags
    }
  }

  const searchQuerySchema = createSearchQuerySchema(
    entitySchema as TObject,
    entityKey
  )

  // POST /entity - Create single entity
  app.post(
    prefix,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        body: createSchema,
        description: `Create a new ${singularize(entityKey)}`,
        operationId: 'create' + capitalize(toCamelCase(singularize(entityKey))),
        response: {
          201: LegacyRef(DocumentSchemaFactory(entitySchema))
          // 400: { $ref: 'error-document' },
          // 401: { $ref: 'unauthorized-response' },
          // 403: { $ref: 'forbidden-response' }
        },
        summary: `Create a new ${singularize(entityKey)}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    //@ts-ignore
    async (req) => {
      return service.create(req.body)
    }
  )

  // GET /entity - Find many entities
  app.get(
    prefix,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Find many ${entityKey}`,
        operationId: 'findMany' + capitalize(toCamelCase(entityKey)),
        querystring: searchQuerySchema,
        response: {
          200: Type.Object({
            data: Type.Array(LegacyRef(entitySchema)),
            meta: Type.Optional(
              Type.Object({
                count: Type.Number(),
                page: Type.Number(),
                pageSize: Type.Number(),
                total: Type.Number()
              })
            )
          })
        },
        summary: `Find many ${entityKey}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    async (request) => {
      return service.findMany(request.query)
    }
  )

  // GET /entity/:id - Find single entity by ID
  app.get(
    `${prefix}/:id`,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Find a${vf(entityKey)} ${singularize(entityKey)}`,
        operationId: 'getOne' + capitalize(toCamelCase(singularize(entityKey))),
        params: Type.Object({
          id: Type.String()
        }),
        response: {
          200: DocumentSchemaFactory(entitySchema),
          404: LegacyRef(NotFoundResponseSchema)
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
  app.patch(
    `${prefix}/:id`,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        body: updateSchema,
        description: `Update a${vf(entityKey)} ${singularize(entityKey)}`,
        operationId: 'update' + capitalize(toCamelCase(singularize(entityKey))),
        params: Type.Object({
          id: Type.String()
        }),
        response: {
          200: DocumentSchemaFactory(entitySchema),
          404: LegacyRef(NotFoundResponseSchema)
        },
        summary: `Update a${vf(entityKey)} ${singularize(entityKey)}`,
        tags: [toTitleCase(entityKey)]
      }
    },
    async (request) => {
      return {
        //@ts-ignore
        data: await service.update(request.params.id, request.body)
      }
    }
  )

  // DELETE /entity/:id - Delete single entity
  app.delete(
    `${prefix}/:id`,
    {
      ...baseRouteOptions,
      schema: {
        ...baseRouteOptions.schema,
        description: `Delete a${vf(entityKey)} ${singularize(entityKey)}`,
        operationId: 'delete' + capitalize(toCamelCase(singularize(entityKey))),
        params: Type.Object({
          id: Type.String()
        }),
        response: {
          200: DocumentSchemaFactory(entitySchema),
          404: LegacyRef(NotFoundResponseSchema)
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

  app.addSchema(entitySchema)
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
