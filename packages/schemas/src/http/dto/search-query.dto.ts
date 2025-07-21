import { z } from 'zod'

import type { BaseEntity } from '@archesai/schemas'

// ==========================================
// FILTER OPERATORS FROM YOUR DATATABLECONFIG
// ==========================================

export const FilterOperation = [
  'eq',
  'ne',
  'lt',
  'lte',
  'gt',
  'gte',
  'iLike',
  'notILike',
  'inArray',
  'notInArray',
  'isEmpty',
  'isNotEmpty',
  'isBetween',
  'isRelativeToToday'
] as const

export type FilterOperationType = (typeof FilterOperation)[number]

export const FilterValueSchema: z.ZodUnion<
  readonly [
    z.ZodString,
    z.ZodNumber,
    z.ZodBoolean,
    z.ZodArray<z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean]>>,
    z.ZodObject<{
      from: z.ZodUnion<readonly [z.ZodString, z.ZodNumber]>
      to: z.ZodUnion<readonly [z.ZodString, z.ZodNumber]>
    }>,
    z.ZodObject<{
      unit: z.ZodEnum<{
        days: 'days'
        months: 'months'
        weeks: 'weeks'
        years: 'years'
      }>
      value: z.ZodNumber
    }>
  ]
> = z
  .union([
    z.string(),
    z.number(),
    z.boolean(),
    z.array(z.union([z.string(), z.number(), z.boolean()])),
    // Range object for isBetween
    z.object({
      from: z.union([z.string(), z.number()]),
      to: z.union([z.string(), z.number()])
    }),
    // Relative date object for isRelativeToToday
    z.object({
      unit: z.enum(['days', 'weeks', 'months', 'years']),
      value: z.number()
    })
  ])
  .meta({
    description:
      'Value for filter conditions, supports strings, numbers, booleans, arrays, ranges, and relative dates',
    id: 'FilterValue'
  })

export type FilterValue = z.infer<typeof FilterValueSchema>

export const OperatorSchema: z.ZodEnum<{
  eq: 'eq'
  gt: 'gt'
  gte: 'gte'
  iLike: 'iLike'
  inArray: 'inArray'
  isBetween: 'isBetween'
  isEmpty: 'isEmpty'
  isNotEmpty: 'isNotEmpty'
  isRelativeToToday: 'isRelativeToToday'
  lt: 'lt'
  lte: 'lte'
  ne: 'ne'
  notILike: 'notILike'
  notInArray: 'notInArray'
}> = z.enum(FilterOperation).meta({
  description: 'Supported filter operations',
  id: 'Operator'
})

// ==========================================
// FILTER CONDITION SCHEMA
// ==========================================

export const FilterConditionSchema: z.ZodObject<{
  field: z.ZodString
  operator: z.ZodEnum<{
    eq: 'eq'
    gt: 'gt'
    gte: 'gte'
    iLike: 'iLike'
    inArray: 'inArray'
    isBetween: 'isBetween'
    isEmpty: 'isEmpty'
    isNotEmpty: 'isNotEmpty'
    isRelativeToToday: 'isRelativeToToday'
    lt: 'lt'
    lte: 'lte'
    ne: 'ne'
    notILike: 'notILike'
    notInArray: 'notInArray'
  }>
  type: z.ZodLiteral<'condition'>
  value: z.ZodUnion<
    readonly [
      z.ZodString,
      z.ZodNumber,
      z.ZodBoolean,
      z.ZodArray<z.ZodUnion<readonly [z.ZodString, z.ZodNumber, z.ZodBoolean]>>,
      z.ZodObject<{
        from: z.ZodUnion<readonly [z.ZodString, z.ZodNumber]>
        to: z.ZodUnion<readonly [z.ZodString, z.ZodNumber]>
      }>,
      z.ZodObject<{
        unit: z.ZodEnum<{
          days: 'days'
          months: 'months'
          weeks: 'weeks'
          years: 'years'
        }>
        value: z.ZodNumber
      }>
    ]
  >
}> = z
  .object({
    field: z.string(),
    operator: OperatorSchema,
    type: z.literal('condition'),
    value: FilterValueSchema
  })
  .meta({
    description: 'A single filter condition with field, operator, and value',
    id: 'FilterCondition'
  })

// ==========================================
// RECURSIVE FILTER NODE SCHEMA
// ==========================================

export type FilterConditionType = z.infer<typeof FilterConditionSchema>

export interface FilterGroupType {
  children: FilterNodeType[]
  operator: 'and' | 'or'
  type: 'group'
}

export type FilterNodeType = FilterConditionType | FilterGroupType

export const FilterNodeSchema: z.ZodType<FilterNodeType> = z
  .discriminatedUnion('type', [
    FilterConditionSchema,
    z
      .object({
        get children() {
          return z.array(FilterNodeSchema)
        },
        operator: z.enum(['and', 'or']),
        type: z.literal('group')
      })
      .describe('A logical group of filter conditions or other groups')
  ])
  .meta({
    description: 'A recursive filter node that can be a condition or group',
    id: 'FilterNode'
  })

// ==========================================
// PAGINATION SCHEMA
// ==========================================

export const PageSchema: z.ZodObject<{
  number: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
  size: z.ZodOptional<z.ZodDefault<z.ZodNumber>>
}> = z
  .object({
    number: z.number().int().min(1).max(Number.MAX_VALUE).default(1).optional(),
    size: z.number().int().min(1).max(100).default(10).optional()
  })
  .meta({
    description: 'Pagination configuration with page number and size',
    id: 'Page'
  })

// ==========================================
// SORT SCHEMA
// ==========================================

export const SortSchema: z.ZodObject<{
  field: z.ZodString
  order: z.ZodEnum<{
    asc: 'asc'
    desc: 'desc'
  }>
}> = z
  .object({
    field: z.string(),
    order: z.enum(['asc', 'desc'])
  })
  .meta({
    description: 'Sorting configuration with field and order',
    id: 'Sort'
  })

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

export interface FilterCondition<TEntity extends BaseEntity> {
  field: keyof TEntity
  operator: FilterOperationType
  type: 'condition'
  value:
    | (boolean | number | string)[]
    | boolean
    | number
    | string
    | { from: number | string; to: number | string }
    | { unit: 'days' | 'months' | 'weeks' | 'years'; value: number }
}

export interface FilterGroup<TEntity extends BaseEntity> {
  children: FilterNode<TEntity>[]
  operator: 'and' | 'or'
  type: 'group'
}

export type FilterNode<TEntity extends BaseEntity> =
  | FilterCondition<TEntity>
  | FilterGroup<TEntity>

export interface SearchQuery<TEntity extends BaseEntity> {
  filter?: FilterNode<TEntity>
  page?: {
    number?: number
    size?: number
  }
  sort?: {
    field: keyof TEntity
    order: 'asc' | 'desc'
  }[]
}

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
