import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { BaseEntity } from '@archesai/domain'

import { LegacyRef } from '@archesai/domain'

import { toTitleCaseNoSpaces } from '#utils/strings'

export const FilterOperation = [
  'equals',
  'gt',
  'gte',
  'in',
  'is_not_null',
  'is_null',
  'like',
  'lt',
  'lte',
  'not_equals',
  'not_in'
] as const
export type FilterOperationType = (typeof FilterOperation)[number]

/**
 * Each filter operation can hold either
 * - a string
 * - a number
 * - a boolean
 * - an array of [string|number|boolean]
 */
export const FilterValue = Type.Union([
  Type.String(),
  Type.Number(),
  Type.Boolean(),
  Type.Array(Type.Union([Type.String(), Type.Number(), Type.Boolean()]))
])

/**
 * A reusable "FieldFilter" schema that maps each operation
 * (like equals, GT, GTE, etc.) to a FilterValue.
 *
 */
export const FieldFilterSchema = Type.Partial(
  Type.Record(
    Type.Union(FilterOperation.map((op) => Type.Literal(op))),
    FilterValue
  ),
  {
    $id: 'FieldFilter',
    description: 'Key-value pairs for filter operations',
    examples: [
      {
        equals: 'value',
        GT: 1,
        IN: ['val1', 'val2']
      }
    ],
    title: 'Field Filter'
  }
)

export const PageSchema = Type.Object(
  {
    number: Type.Optional(
      Type.Integer({ default: 1, maximum: Number.MAX_VALUE, minimum: 0 })
    ),
    size: Type.Optional(Type.Integer({ default: 10, maximum: 100, minimum: 1 }))
  },
  {
    description: 'Pagination',
    examples: [
      {
        number: 1,
        size: 10
      }
    ],
    title: 'Page'
  }
)

export const SortSchema = Type.String({
  $id: 'Sort',
  description: 'Sort by name ascending and createdAt descending',
  examples: ['name,-createdAt'],
  title: 'Sort'
})

export const SearchQuerySchema = Type.Object(
  {
    page: Type.Optional(PageSchema),
    sort: Type.Optional(SortSchema)
  },
  {
    $id: 'SearchQuery',
    description: 'Search query',
    title: 'Search Query'
  }
)

const createFilterSchema = (EntitySchema: TObject) => {
  const entityFields = Type.Union(
    Object.keys(EntitySchema.properties).map((key) => Type.Literal(key))
  )

  return Type.Record(entityFields, LegacyRef(FieldFilterSchema), {
    // $id: 'Filter',
    description: 'Filter',
    title: 'Filter'
  })
}

export const createSearchQuerySchema = (
  EntitySchema: TObject,
  entityKey: string
) => {
  return Type.Composite(
    [
      Type.Object({
        filter: Type.Optional(createFilterSchema(EntitySchema))
      }),
      SearchQuerySchema
    ],
    {
      // $id: `Search${toTitleCaseNoSpaces(entityKey)}Query`,
      description: `Search ${toTitleCaseNoSpaces(entityKey)} query`,
      title: `Search ${toTitleCaseNoSpaces(entityKey)} Query`
    }
  )
}

export type EntityFilter<TEntity extends BaseEntity> = Partial<
  Record<
    keyof TEntity,
    Partial<
      Record<
        FilterOperationType,
        (boolean | number | string)[] | boolean | number | string
      >
    >
  >
>

export interface SearchQuery<TEntity extends BaseEntity> {
  filter?: EntityFilter<TEntity>
  page?: {
    number?: number
    size?: number
  }
  sort?: string
}
