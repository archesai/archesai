import type { TObject } from '@sinclair/typebox'
import type { FastifyPluginCallback } from 'fastify'

import { Type } from '@sinclair/typebox'

import type { BaseEntity, BaseInsertion } from '@archesai/schemas'

import { LegacyRef } from '@archesai/schemas'

import type { BaseService } from '#common/base.service'
import type { Controller } from '#http/interfaces/controller.interface'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

import { ArchesApiNotFoundResponseSchema } from '#exceptions/schemas/arches-api-not-found-response.schema'
import { createSearchQuerySchema } from '#http/dto/search-query.dto'
import { AuthenticatedGuard } from '#http/guards/authenticated.guard'
import { DocumentSchemaFactory } from '#http/schemas/document.schema'
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
  public readonly entityKey: string
  public readonly [IS_CONTROLLER] = true
  protected readonly createSchema: TObject
  protected readonly service: BaseService<TEntity, TInsert>
  protected readonly updateSchema: TObject
  private readonly entitySchema: TObject
  // private readonly logger = new Logger(BaseController.name)
  private readonly searchQuerySchema: TObject

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

    this.searchQuerySchema = createSearchQuerySchema(
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
            201: DocumentSchemaFactory(this.entitySchema)
          },
          summary: `Create a new ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        return {
          data: await this.service.create(request.body as TCreateRequest)
        }
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
            200: DocumentSchemaFactory(this.entitySchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Delete a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        return {
          data: await this.service.delete(request.params.id)
        }
      }
    )

    app.get(
      `/${this.entityKey}`,
      {
        preValidation: [AuthenticatedGuard()],
        schema: {
          description: `Find many ${this.entityKey}`,
          operationId: 'findMany' + capitalize(toCamelCase(this.entityKey)),
          querystring: this.searchQuerySchema,
          response: {
            200: Type.Object({
              data: Type.Array(LegacyRef(this.entitySchema)),
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
          summary: `Find many ${this.entityKey}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        return this.service.findMany(request.query)
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
            200: DocumentSchemaFactory(this.entitySchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Find a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        return {
          data: await this.service.findOne(request.params.id)
        }
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
            200: DocumentSchemaFactory(this.entitySchema),
            404: LegacyRef(ArchesApiNotFoundResponseSchema)
          },
          summary: `Update a${vf(this.entityKey)} ${singularize(this.entityKey)}`,
          tags: [toTitleCase(this.entityKey)]
        }
      },
      async (request) => {
        return {
          data: await this.service.update(
            request.params.id,
            request.body as TUpdateRequest
          )
        }
      }
    )

    app.addSchema(this.entitySchema)
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
}
