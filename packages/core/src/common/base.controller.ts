import type { StaticDecode, TObject } from '@sinclair/typebox'
import type { FastifyReply, FastifyRequest } from 'fastify'

import { Type } from '@sinclair/typebox'

import type { BaseEntity, BaseInsertion } from '@archesai/domain'

import { LegacyRef } from '@archesai/domain'

import type { BaseService } from '#common/base.service'
import type { SearchQuery } from '#http/dto/search-query.dto'
import type { Controller } from '#http/interfaces/controller.interface'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

import { ArchesApiNotFoundResponseSchema } from '#exceptions/schemas/arches-api-not-found-response.schema'
import { createSearchQuerySchema } from '#http/dto/search-query.dto'
import { createCollectionResponseSchema } from '#http/factories/collection-response.schema'
import { createIndividualResponseSchema } from '#http/factories/individual-response.schema'
import { createResourceObjectSchema } from '#http/factories/resource-object.schema'
import { AuthenticatedGuard } from '#http/guards/authenticated.guard'
import { capitalize } from '#utils/capitalize'
import { singularize } from '#utils/pluralize'
import { toCamelCase, toTitleCase, vf } from '#utils/strings'

export const IS_CONTROLLER = Symbol.for('isController')

/**
 * A base controller for handling CRUD operations on a resource.
 */
export abstract class BaseController<
  TEntity extends BaseEntity = BaseEntity,
  TInsert extends BaseInsertion<TEntity> = BaseInsertion<TEntity>,
  TCreateRequest extends TInsert = TInsert,
  TUpdateRequest extends Partial<TInsert> = Partial<TInsert>
> implements Controller
{
  public readonly [IS_CONTROLLER] = true
  protected readonly CreateSchema: TObject
  protected readonly IndividualEntityResponseSchema: TObject
  protected readonly PaginatedEntityResponseSchema: TObject
  protected readonly service: BaseService<TEntity, TInsert>
  protected readonly UpdateSchema: TObject
  private readonly entityKey: string
  private readonly EntitySchema: TObject
  private readonly ResourceObjectSchema: TObject
  private readonly SearchQuerySchema: TObject

  constructor(
    entityKey: string,
    EntitySchema: TObject,
    CreateSchema: TObject,
    UpdateSchema: TObject,
    service: BaseService<TEntity, TInsert>
  ) {
    this.entityKey = entityKey
    this.EntitySchema = EntitySchema
    this.service = service
    this.CreateSchema = CreateSchema
    this.UpdateSchema = UpdateSchema
    this.ResourceObjectSchema = createResourceObjectSchema(
      this.EntitySchema,
      this.entityKey
    )
    this.IndividualEntityResponseSchema = createIndividualResponseSchema(
      this.ResourceObjectSchema,
      this.entityKey
    )
    this.PaginatedEntityResponseSchema = createCollectionResponseSchema(
      this.ResourceObjectSchema,
      this.entityKey
    )
    this.SearchQuerySchema = createSearchQuerySchema(
      this.EntitySchema,
      this.entityKey
    )
  }

  public async create(
    request: FastifyRequest<{
      Body: TCreateRequest
      Params: { id: string }
    }>,
    _reply: FastifyReply
  ): Promise<StaticDecode<typeof this.IndividualEntityResponseSchema>> {
    const created = await this.service.create(request.body as TCreateRequest)
    return this.toIndividualResponse(created)
  }

  public async delete(
    request: FastifyRequest<{
      Params: { id: string }
    }>,
    _reply: FastifyReply
  ): Promise<StaticDecode<typeof this.IndividualEntityResponseSchema>> {
    const deleted = await this.service.delete(request.params.id)
    return this.toIndividualResponse(deleted)
  }

  public async findMany(
    request: FastifyRequest<{
      Querystring: SearchQuery<TEntity>
    }>,
    _reply: FastifyReply
    // query: Static<typeof this.SearchQuerySchema>
  ): Promise<StaticDecode<typeof this.PaginatedEntityResponseSchema>> {
    const found = await this.service.findMany(request.query)
    return this.toPaginatedResponse({
      count: found.count,
      data: found.data,
      query: request.query
    })
  }

  public async findOne(
    request: FastifyRequest<{
      Params: { id: string }
    }>,
    _reply: FastifyReply
  ): Promise<StaticDecode<typeof this.IndividualEntityResponseSchema>> {
    const found = await this.service.findOne(request.params.id)
    return this.toIndividualResponse(found)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/${this.entityKey}`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          body: this.CreateSchema,
          description: `Create a new ${singularize(this.entityKey)}`,
          operationId:
            'create' + capitalize(toCamelCase(singularize(this.entityKey))),
          response: {
            201: this.IndividualEntityResponseSchema
          },
          summary: `Create a new ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.create.bind(this)
    )

    app.delete(
      `/${this.entityKey}/:id`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `Delete a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          operationId:
            'delete' + capitalize(toCamelCase(singularize(this.entityKey))),
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: this.IndividualEntityResponseSchema,
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Delete a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.delete.bind(this)
    )

    app.get(
      `/${this.entityKey}`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `Find many ${this.entityKey}`,
          operationId: 'findMany' + capitalize(toCamelCase(this.entityKey)),
          querystring: this.SearchQuerySchema,
          response: {
            200: this.PaginatedEntityResponseSchema
          },
          summary: `Find many ${this.entityKey}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.findMany.bind(this)
    )

    app.get(
      `/${this.entityKey}/:id`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `Find a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          operationId:
            'getOne' + capitalize(toCamelCase(singularize(this.entityKey))),
          params: {
            properties: {
              id: { type: 'string' }
            },
            required: ['id'],
            type: 'object'
          },
          response: {
            200: this.IndividualEntityResponseSchema,
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Find a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.findOne.bind(this)
    )

    app.patch(
      `/${this.entityKey}/:id`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          body: this.UpdateSchema,
          description: `Update a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          operationId:
            'update' + capitalize(toCamelCase(singularize(this.entityKey))),
          params: {
            properties: {
              id: { type: 'string' }
            },
            required: ['id'],
            type: 'object'
          },
          response: {
            200: this.IndividualEntityResponseSchema,
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Update a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      this.update.bind(this)
    )
  }

  public async update(
    request: FastifyRequest<{
      Body: TUpdateRequest
      Params: { id: string }
    }>,
    _reply: FastifyReply
  ): Promise<StaticDecode<typeof this.IndividualEntityResponseSchema>> {
    const updated = await this.service.update(
      request.params.id,
      request.body as TUpdateRequest
    )
    return this.toIndividualResponse(updated)
  }

  protected toIndividualResponse(
    input: TEntity
  ): StaticDecode<typeof this.IndividualEntityResponseSchema> {
    const { id, type, ...attributes } = input
    return {
      data: {
        attributes: attributes,
        id: id,
        type: type
      },
      links: {
        self: `${type}s/${input.id}`
      }
    }
  }

  protected toPaginatedResponse(input: {
    count: number
    data: TEntity[]
    query?: SearchQuery<TEntity>
  }): StaticDecode<typeof this.PaginatedEntityResponseSchema> {
    return {
      data: input.data.map((entity) => {
        const { id, type, ...attributes } = entity
        return {
          attributes: attributes,
          id: id,
          relationships: {},
          type: type
        }
      }),
      links: {}
      // meta: {
      //   // query: query,
      //   total_pages: 1,
      //   total_records: input.count
      // }
    }
  }
}

// Post('bulk')
// async createMany(
//   Body() createManyRequest: TCreateRequest[]
// ): Promise<PaginatedJsonApiResponse<TEntity>> {
//   const created = await this.service.createMany(createManyRequest)
//   return this.toPaginatedResponse({
//     count: created.count,
//     data: created.data,
//     type: ResponseObject.name
//   })
// }

// Patch('bulk')
// async updateMany(
//   Body() updateManyRequest: TUpdateRequest,
//   Query() query: SearchQuery<TEntity>
// ): Promise<PaginatedJsonApiResponse<TEntity>> {
//   const updated = await this.service.updateMany(updateManyRequest, query)
//   return this.toPaginatedResponse({
//     count: updated.count,
//     data: updated.data,
//     query,
//     type: ResponseObject.name
//   })
// }

// Post('bulk')
// async deleteMany(
//   Query() query: SearchQuery<TEntity>
// ): Promise<PaginatedJsonApiResponse<TEntity>> {
//   const deleted = await this.service.deleteMany(query)
//   return this.toPaginatedResponse({
//     count: deleted.count,
//     data: deleted.data,
//     query,
//     type: ResponseObject.name
//   })
// }
