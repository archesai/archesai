import { z } from 'zod'

import type { FilterGroupType, FilterNodeType } from '#query/filter-node.schema'

import { FilterNodeSchema } from '#query/filter-node.schema'
import { FilterValueSchema } from '#query/filter-value.schema'
import { OperatorSchema } from '#query/operator.schema'
import { PageSchema } from '#query/page.schema'
import { SortSchema } from '#query/sort.schema'

// ==========================================
// MAIN SEARCH QUERY SCHEMA
// ==========================================

export const SearchQuerySchema: z.ZodObject<{
  filter: z.ZodOptional<
    z.ZodType<FilterNodeType, unknown, z.core.$ZodTypeInternals<FilterNodeType>>
  >
  page: z.ZodOptional<
    z.ZodObject<{
      number: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
      size: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
    }>
  >
  sort: z.ZodOptional<
    z.ZodArray<
      z.ZodObject<{
        field: z.ZodString
        order: z.ZodEnum<{
          asc: 'asc'
          desc: 'desc'
        }>
      }>
    >
  >
}> = z
  .object({
    filter: FilterNodeSchema.optional(),
    page: PageSchema.optional(),
    sort: z.array(SortSchema).optional()
  })
  .meta({
    description:
      'Complete search query with nested filters, pagination, and sorting',
    id: 'SearchQuery'
  })

// ==========================================
// TYPESCRIPT INTERFACES
// ==========================================

// ==========================================
// ENTITY-SPECIFIC SCHEMA CREATION
// ==========================================

export const createSearchQuerySchema = (
  entitySchema: z.ZodObject,
  entityKey: string
): z.ZodObject<{
  filter: z.ZodOptional<
    z.ZodType<
      | FilterGroupType
      | {
          field: string
          operator:
            | 'eq'
            | 'gt'
            | 'gte'
            | 'iLike'
            | 'inArray'
            | 'isBetween'
            | 'isEmpty'
            | 'isNotEmpty'
            | 'isRelativeToToday'
            | 'lt'
            | 'lte'
            | 'ne'
            | 'notILike'
            | 'notInArray'
          type: 'condition'
          value:
            | (boolean | number | string)[]
            | boolean
            | number
            | string
            | {
                from: number | string
                to: number | string
              }
            | {
                unit: 'days' | 'months' | 'weeks' | 'years'
                value: number
              }
        },
      unknown,
      z.core.$ZodTypeInternals<
        | FilterGroupType
        | {
            field: string
            operator:
              | 'eq'
              | 'gt'
              | 'gte'
              | 'iLike'
              | 'inArray'
              | 'isBetween'
              | 'isEmpty'
              | 'isNotEmpty'
              | 'isRelativeToToday'
              | 'lt'
              | 'lte'
              | 'ne'
              | 'notILike'
              | 'notInArray'
            type: 'condition'
            value:
              | (boolean | number | string)[]
              | boolean
              | number
              | string
              | {
                  from: number | string
                  to: number | string
                }
              | {
                  unit: 'days' | 'months' | 'weeks' | 'years'
                  value: number
                }
          }
      >
    >
  >
  page: z.ZodOptional<
    z.ZodObject<{
      number: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
      size: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
    }>
  >
  sort: z.ZodOptional<
    z.ZodArray<
      z.ZodObject<{
        field: z.ZodEnum<Record<string, string>>
        order: z.ZodEnum<{
          asc: 'asc'
          desc: 'desc'
        }>
      }>
    >
  >
}> => {
  // Extract field names from Zod schema
  const entityFields = Object.keys(entitySchema.shape)

  // Create entity-specific filter condition
  const EntityFilterConditionSchema = z.object({
    field: z.enum(entityFields),
    operator: OperatorSchema,
    type: z.literal('condition'),
    value: FilterValueSchema
  })

  // Create entity-specific filter node type
  type EntityFilterConditionType = z.infer<typeof EntityFilterConditionSchema>

  interface EntityFilterGroupType {
    children: EntityFilterNodeType[]
    operator: 'and' | 'or'
    type: 'group'
  }

  type EntityFilterNodeType = EntityFilterConditionType | EntityFilterGroupType

  // Create entity-specific recursive filter node
  const EntityFilterNodeSchema: z.ZodType<EntityFilterNodeType> = z
    .discriminatedUnion('type', [
      EntityFilterConditionSchema,
      z.object({
        get children() {
          return z.array(EntityFilterNodeSchema)
        },
        operator: z.enum(['and', 'or']),
        type: z.literal('group')
      })
    ])
    .meta({
      description: `A recursive filter node for ${entityKey} entity that can be a condition or group`,
      id: `${entityKey}FilterNode`
    })

  // Create entity-specific sort schema
  const EntitySortSchema = z
    .object({
      field: z.enum(entityFields),
      order: z.enum(['asc', 'desc'])
    })
    .meta({
      description: `Sorting configuration for ${entityKey} entity with field and order`,
      id: `${entityKey}Sort`
    })

  const searchQuerySchema = z
    .object({
      filter: EntityFilterNodeSchema.optional(),
      page: PageSchema.optional(),
      sort: z.array(EntitySortSchema).optional()
    })
    .meta({
      description: `Search query schema for ${entityKey} entity`
      // id: `${entityKey}SearchQuery`
    })

  return searchQuerySchema
}
