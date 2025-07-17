/* eslint-disable @typescript-eslint/no-unnecessary-condition */
import type {
  TArray,
  TBoolean,
  TInteger,
  TLiteral,
  TNumber,
  TObject,
  TOptional,
  TRecursive,
  TString,
  TThis,
  TUnion,
  TUnsafe
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { BaseEntity } from '@archesai/schemas'

import { LegacyRef } from '@archesai/schemas'

function toTitleCaseNoSpaces(str: string): string {
  return (
    str
      // Replace underscores/dashes with spaces
      .replace(/[_-]+/g, ' ')
      .trim()
      .split(/\s+/)
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join('')
  )
}

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

// ==========================================
// FILTER VALUES WITH FULL SUPPORT
// ==========================================

export const FilterValueSchema: TUnion<
  [
    TString,
    TNumber,
    TBoolean,
    TArray<TUnion<[TString, TNumber, TBoolean]>>,
    TObject<{
      from: TUnion<[TString, TNumber]>
      to: TUnion<[TString, TNumber]>
    }>,
    TObject<{
      unit: TUnion<
        [
          TLiteral<'days'>,
          TLiteral<'weeks'>,
          TLiteral<'months'>,
          TLiteral<'years'>
        ]
      >
      value: TNumber
    }>
  ]
> = Type.Union(
  [
    Type.String(),
    Type.Number(),
    Type.Boolean(),
    Type.Array(Type.Union([Type.String(), Type.Number(), Type.Boolean()])),
    // Range object for isBetween
    Type.Object({
      from: Type.Union([Type.String(), Type.Number()]),
      to: Type.Union([Type.String(), Type.Number()])
    }),
    // Relative date object for isRelativeToToday
    Type.Object({
      unit: Type.Union([
        Type.Literal('days'),
        Type.Literal('weeks'),
        Type.Literal('months'),
        Type.Literal('years')
      ]),
      value: Type.Number()
    })
  ],
  {
    $id: 'FilterValueSchema',
    description:
      'Value for filter conditions, supports strings, numbers, booleans, arrays, ranges, and relative dates',
    examples: [
      'John%',
      25,
      true,
      ['apple', 'banana'],
      { from: 10, to: 20 },
      { unit: 'days', value: 7 }
    ],
    title: 'Filter Value'
  }
)

export const OperatorSchema: TUnion<
  TLiteral<
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
  >[]
> = Type.Union(
  FilterOperation.map((op) => Type.Literal(op)),
  {
    $id: 'Operator',
    description: 'Supported filter operators',
    examples: FilterOperation
  }
)

// ==========================================
// FILTER CONDITION SCHEMA
// ==========================================

export const FilterConditionSchema: TObject<{
  field: TString
  operator: TUnsafe<
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
  >
  type: TLiteral<'condition'>
  value: TUnsafe<
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
  >
}> = Type.Object(
  {
    field: Type.String(),
    operator: OperatorSchema,
    type: Type.Literal('condition'),
    value: FilterValueSchema
  },
  {
    $id: 'FilterCondition',
    description: 'A single filter condition with field, operator, and value',
    examples: [
      {
        field: 'name',
        operator: 'iLike',
        type: 'condition',
        value: 'John%'
      },
      {
        field: 'age',
        operator: 'isBetween',
        type: 'condition',
        value: { from: 25, to: 35 }
      }
    ],
    title: 'Filter Condition'
  }
)

// ==========================================
// RECURSIVE FILTER NODE SCHEMA
// ==========================================

export const FilterNodeSchema: TRecursive<
  TUnion<
    [
      TUnsafe<{
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
      }>,
      TObject<{
        children: TArray<TThis>
        operator: TUnion<[TLiteral<'and'>, TLiteral<'or'>]>
        type: TLiteral<'group'>
      }>
    ]
  >
> = Type.Recursive(
  (This) =>
    Type.Union([
      FilterConditionSchema,
      Type.Object(
        {
          children: Type.Array(This),
          operator: Type.Union([Type.Literal('and'), Type.Literal('or')]),
          type: Type.Literal('group')
        },
        {
          $id: 'FilterGroup',
          description: 'A logical group of filter conditions or other groups',
          examples: [
            {
              children: [
                {
                  field: 'name',
                  operator: 'iLike',
                  type: 'condition',
                  value: 'John%'
                },
                {
                  field: 'age',
                  operator: 'gt',
                  type: 'condition',
                  value: 25
                }
              ],
              operator: 'and',
              type: 'group'
            }
          ],
          title: 'Filter Group'
        }
      )
    ]),
  {
    $id: 'FilterNode',
    description: 'A recursive filter node that can be a condition or group',
    title: 'Filter Node'
  }
)

// ==========================================
// PAGINATION SCHEMA
// ==========================================

export const PageSchema: TObject<{
  number: TOptional<TInteger>
  size: TOptional<TInteger>
}> = Type.Object(
  {
    number: Type.Optional(
      Type.Integer({ default: 1, maximum: Number.MAX_VALUE, minimum: 1 })
    ),
    size: Type.Optional(Type.Integer({ default: 10, maximum: 100, minimum: 1 }))
  },
  {
    $id: 'Page',
    description: 'Pagination configuration',
    examples: [
      {
        number: 1,
        size: 10
      }
    ],
    title: 'Page'
  }
)

// ==========================================
// SORT SCHEMA
// ==========================================

export const SortSchema: TObject<{
  field: TString
  order: TUnion<[TLiteral<'asc'>, TLiteral<'desc'>]>
}> = Type.Object(
  {
    field: Type.String(),
    order: Type.Union([Type.Literal('asc'), Type.Literal('desc')])
  },
  {
    $id: 'Sort',
    description: 'Sort configuration',
    examples: [
      { field: 'name', order: 'asc' },
      { field: 'createdAt', order: 'desc' }
    ],
    title: 'Sort'
  }
)

// ==========================================
// MAIN SEARCH QUERY SCHEMA
// ==========================================

export const SearchQuerySchema: TObject<{
  filter: TOptional<
    TRecursive<
      TUnion<
        [
          TObject<{
            field: TString
            operator: TUnion<
              TLiteral<
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
              >[]
            >
            type: TLiteral<'condition'>
            value: TUnion<
              [
                TString,
                TNumber,
                TBoolean,
                TArray<TUnion<[TString, TNumber, TBoolean]>>,
                TObject<{
                  from: TUnion<[TString, TNumber]>
                  to: TUnion<[TString, TNumber]>
                }>,
                TObject<{
                  unit: TUnion<
                    [
                      TLiteral<'days'>,
                      TLiteral<'weeks'>,
                      TLiteral<'months'>,
                      TLiteral<'years'>
                    ]
                  >
                  value: TNumber
                }>
              ]
            >
          }>,
          TObject<{
            children: TArray<TThis>
            operator: TUnion<[TLiteral<'and'>, TLiteral<'or'>]>
            type: TLiteral<'group'>
          }>
        ]
      >
    >
  >
  page: TOptional<
    TObject<{
      number: TOptional<TInteger>
      size: TOptional<TInteger>
    }>
  >
  sort: TOptional<
    TArray<
      TObject<{
        field: TString
        order: TUnion<[TLiteral<'asc'>, TLiteral<'desc'>]>
      }>
    >
  >
}> = Type.Object(
  {
    filter: Type.Optional(FilterNodeSchema),
    page: Type.Optional(PageSchema),
    sort: Type.Optional(Type.Array(SortSchema))
  },
  {
    $id: 'SearchQuery',
    description:
      'Complete search query with nested filters, pagination, and sorting',
    examples: [
      {
        filter: {
          children: [
            {
              field: 'name',
              operator: 'iLike',
              type: 'condition',
              value: 'John%'
            },
            {
              children: [
                {
                  field: 'age',
                  operator: 'gt',
                  type: 'condition',
                  value: 25
                },
                {
                  field: 'department',
                  operator: 'eq',
                  type: 'condition',
                  value: 'Engineering'
                }
              ],
              operator: 'or',
              type: 'group'
            }
          ],
          operator: 'and',
          type: 'group'
        },
        page: { number: 1, size: 10 },
        sort: [{ field: 'name', order: 'asc' }]
      }
    ],
    title: 'Search Query'
  }
)

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
  entitySchema: TObject,
  entityKey: string
): TObject<{
  filter: TOptional<
    TRecursive<
      TUnion<
        [
          TObject<{
            field: TUnion<TLiteral<string>[]>
            operator: TUnsafe<
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
            >
            type: TLiteral<'condition'>
            value: TUnsafe<
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
            >
          }>,
          TObject<{
            children: TArray<TThis>
            operator: TUnion<[TLiteral<'and'>, TLiteral<'or'>]>
            type: TLiteral<'group'>
          }>
        ]
      >
    >
  >
  page: TOptional<
    TUnsafe<{
      number?: number | undefined
      size?: number | undefined
    }>
  >
  sort: TOptional<
    TArray<
      TObject<{
        field: TUnion<TLiteral<string>[]>
        order: TUnion<[TLiteral<'asc'>, TLiteral<'desc'>]>
      }>
    >
  >
}> => {
  // Create entity-specific field validation
  const entityFields = Object.keys(entitySchema.properties)

  // Create entity-specific filter condition
  const EntityFilterConditionSchema = Type.Object(
    {
      field: Type.Union(entityFields.map((field) => Type.Literal(field))),
      operator: LegacyRef(OperatorSchema),
      type: Type.Literal('condition'),
      value: LegacyRef(FilterValueSchema)
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}FilterCondition`,
      title: `${toTitleCaseNoSpaces(entityKey)} Filter Condition`
    }
  )

  // Create entity-specific recursive filter node
  const EntityFilterNodeSchema = Type.Recursive(
    (This) =>
      Type.Union([
        EntityFilterConditionSchema,
        Type.Object({
          children: Type.Array(This),
          operator: Type.Union([Type.Literal('and'), Type.Literal('or')]),
          type: Type.Literal('group')
        })
      ]),
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}FilterNode`,
      title: `${toTitleCaseNoSpaces(entityKey)} Filter Node`
    }
  )

  // Create entity-specific sort schema
  const EntitySortSchema = Type.Object(
    {
      field: Type.Union(entityFields.map((field) => Type.Literal(field))),
      order: Type.Union([Type.Literal('asc'), Type.Literal('desc')])
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}Sort`,
      title: `${toTitleCaseNoSpaces(entityKey)} Sort`
    }
  )

  return Type.Object(
    {
      filter: Type.Optional(EntityFilterNodeSchema),
      page: Type.Optional(LegacyRef(PageSchema)),
      sort: Type.Optional(Type.Array(EntitySortSchema))
    },
    {
      $id: `Search${toTitleCaseNoSpaces(entityKey)}Query`,
      description: `Complete search query for ${entityKey} with nested filters, pagination, and sorting`,
      title: `Search ${toTitleCaseNoSpaces(entityKey)} Query`
    }
  )
}

// ==========================================
// UTILITY FUNCTIONS
// ==========================================

// Count total conditions in a filter tree
export function countConditions<TEntity extends BaseEntity>(
  filter: FilterNode<TEntity>
): number {
  if (filter.type === 'condition') {
    return 1
  } else {
    return filter.children.reduce(
      (sum, child) => sum + countConditions(child),
      0
    )
  }
}

// Convert filter tree to SQL-like string (for debugging)
export function filterToSqlString<TEntity extends BaseEntity>(
  filter: FilterNode<TEntity>
): string {
  if (filter.type === 'condition') {
    const { field, operator, value } = filter
    return `${String(field)} ${operator} ${JSON.stringify(value)}`
  } else {
    const childStrings = filter.children.map((child) =>
      filterToSqlString(child)
    )
    const joinedChildren = childStrings.join(
      ` ${filter.operator.toUpperCase()} `
    )
    return `(${joinedChildren})`
  }
}

// Get all unique fields used in a filter tree
export function getUsedFields<TEntity extends BaseEntity>(
  filter: FilterNode<TEntity>
): (keyof TEntity)[] {
  if (filter.type === 'condition') {
    return [filter.field]
  } else {
    const fields = filter.children.flatMap((child) => getUsedFields(child))
    return [...new Set(fields)]
  }
}

// Validate filter tree structure
export function validateFilterTree<TEntity extends BaseEntity>(
  filter: FilterNode<TEntity>
): { errors: string[]; valid: boolean } {
  const errors: string[] = []

  function validate(node: FilterNode<TEntity>, depth = 0): void {
    if (depth > 10) {
      errors.push('Filter tree too deep (max 10 levels)')
      return
    }

    if (node.type === 'condition') {
      if (!node.field) errors.push('Condition missing field')
      if (!node.operator) errors.push('Condition missing operator')
      if (node.value === undefined) errors.push('Condition missing value')
    } else {
      if (!node.children || node.children.length === 0) {
        errors.push('Group must have at least one child')
      }
      node.children.forEach((child) => {
        validate(child, depth + 1)
      })
    }
  }

  validate(filter)
  return { errors, valid: errors.length === 0 }
}

// ==========================================
// BACKEND MAPPING (if needed for legacy support)
// ==========================================

export const operatorMapping = {
  eq: 'equals',
  gt: 'gt',
  gte: 'gte',
  iLike: 'like',
  inArray: 'in',
  isBetween: 'between',
  isEmpty: 'is_null',
  isNotEmpty: 'is_not_null',
  isRelativeToToday: 'relative_date',
  lt: 'lt',
  lte: 'lte',
  ne: 'not_equals',
  notILike: 'not_like',
  notInArray: 'not_in'
} as const

export type BackendOperator =
  (typeof operatorMapping)[keyof typeof operatorMapping]

// Convert frontend filter to backend format
export function convertToBackendFormat<TEntity extends BaseEntity>(
  filter: FilterNode<TEntity>
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): any {
  if (filter.type === 'condition') {
    return {
      field: filter.field,
      operator: operatorMapping[filter.operator] || filter.operator,
      value: filter.value
    }
  } else {
    return {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-return
      conditions: filter.children.map((child) => convertToBackendFormat(child)),
      operator: filter.operator
    }
  }
}
