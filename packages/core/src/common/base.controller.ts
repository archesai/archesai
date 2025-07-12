import type { StaticDecode, TObject } from '@sinclair/typebox'
import type { FastifyPluginCallback } from 'fastify'

import { Type } from '@sinclair/typebox'

import type { BaseEntity, BaseInsertion } from '@archesai/schemas'

import { LegacyRef } from '@archesai/schemas'

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
  protected readonly createSchema: TObject
  protected readonly IndividualEntityResponseSchema: TObject
  protected readonly PaginatedEntityResponseSchema: TObject
  protected readonly service: BaseService<TEntity, TInsert>
  protected readonly updateSchema: TObject
  private readonly entityKey: string
  private readonly entitySchema: TObject
  private readonly ResourceObjectSchema: TObject
  private readonly SearchQuerySchema: TObject

  constructor(
    entityKey: string,
    entitySchema: TObject,
    createSchema: TObject,
    updateSchema: TObject,
    service: BaseService<TEntity, TInsert>
  ) {
    this.entityKey = entityKey
    this.entitySchema = entitySchema
    this.service = service
    this.createSchema = createSchema
    this.updateSchema = updateSchema
    this.ResourceObjectSchema = createResourceObjectSchema(
      this.entitySchema,
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
      this.entitySchema,
      this.entityKey
    )
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/${this.entityKey}`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          body: this.createSchema,
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
      async (request) => {
        const created = await this.service.create(
          request.body as TCreateRequest
        )
        return this.toIndividualResponse(created)
      }
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
      async (request) => {
        const deleted = await this.service.delete(request.params.id)
        return this.toIndividualResponse(deleted)
      }
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
      async (request) => {
        const found = await this.service.findMany(request.query)
        return this.toPaginatedResponse({
          count: found.count,
          data: found.data,
          query: request.query
        })
      }
    )

    app.get(
      `/${this.entityKey}/:id`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `Find a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          operationId:
            'getOne' + capitalize(toCamelCase(singularize(this.entityKey))),
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: this.IndividualEntityResponseSchema,
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Find a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        const found = await this.service.findOne(request.params.id)
        return this.toIndividualResponse(found)
      }
    )

    app.patch(
      `/${this.entityKey}/:id`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          body: this.updateSchema,
          description: `Update a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          operationId:
            'update' + capitalize(toCamelCase(singularize(this.entityKey))),
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: this.IndividualEntityResponseSchema,
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Update a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        const updated = await this.service.update(
          request.params.id,
          request.body as TUpdateRequest
        )
        return this.toIndividualResponse(updated)
      }
    )
  }

  /**
   * Convert this controller to a Fastify plugin
   */
  public toPlugin(): FastifyPluginCallback {
    return (app: HttpInstance, _, done) => {
      this.registerRoutes(app)
      done()
    }
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
    }
  }
}
