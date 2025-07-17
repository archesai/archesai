import type { Row, RowData } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/schemas'

// Extend TanStack Table column meta for filter configuration
declare module '@tanstack/react-table' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface ColumnMeta<TData extends RowData, TValue> {
    filterOperators?: string[]
    // Filter configuration
    filterVariant:
      | 'boolean'
      | 'date'
      | 'dateRange'
      | 'multiSelect'
      | 'number'
      | 'range'
      | 'select'
      | 'text'

    icon: React.FC<React.SVGProps<SVGSVGElement>>
    // Display options
    label?: string

    max?: number

    // Number/Date range
    min?: number
    // Select/MultiSelect options
    options?: {
      count?: number
      icon?: React.FC<React.SVGProps<SVGSVGElement>>
      label: string
      value: string
    }[]
    step?: number
    unit?: string
  }
}

// Row action interface for data tables
export interface DataTableRowAction<TData> {
  icon?: React.FC<React.SVGProps<SVGSVGElement>>
  label?: string
  onClick?: (row: Row<TData>) => void
  row: Row<TData>
  variant: 'custom' | 'delete' | 'update' | 'view'
}

// Filter operators mapped to your SearchQuery DTO
export const FILTER_OPERATORS = {
  boolean: ['eq', 'ne'],
  date: [
    'eq',
    'ne',
    'lt',
    'lte',
    'gt',
    'gte',
    'isBetween',
    'isRelativeToToday'
  ],
  multiSelect: ['inArray', 'notInArray'],
  number: ['eq', 'ne', 'lt', 'lte', 'gt', 'gte', 'isBetween'],
  select: ['eq', 'ne', 'inArray', 'notInArray'],
  text: ['iLike', 'notILike', 'eq', 'ne', 'isEmpty', 'isNotEmpty']
} as const

// Utility type for extracting column keys
export type ColumnKey<TEntity extends BaseEntity> = Extract<
  keyof TEntity,
  string
>
// Data table configuration
export interface DataTableConfig {
  allowNestedFilters: boolean
  // Pagination
  defaultPageSize: number

  enableBulkActions: boolean
  // Export
  enableExport: boolean

  // Selection
  enableSelection: boolean
  exportFormats: string[]

  // Filters
  maxFilters: number
  pageSizeOptions: number[]
}

export type FilterOperator = (typeof FILTER_OPERATORS)[FilterVariant][number]

export type FilterVariant = keyof typeof FILTER_OPERATORS

// Simple filter condition interface (matches your DTO)
export interface SimpleFilterCondition<TEntity extends BaseEntity> {
  field: keyof TEntity
  operator: FilterOperator
  value:
    | (boolean | number | string)[]
    | boolean
    | number
    | string
    | { from: number | string; to: number | string }
    | { unit: 'days' | 'months' | 'weeks' | 'years'; value: number }
}

export const DEFAULT_DATA_TABLE_CONFIG: DataTableConfig = {
  allowNestedFilters: true,
  defaultPageSize: 10,
  enableBulkActions: true,
  enableExport: true,
  enableSelection: true,
  exportFormats: ['csv', 'json'],
  maxFilters: 10,
  pageSizeOptions: [10, 20, 50, 100]
}

// Helper functions for filter operations
export function getDefaultFilterOperator(
  filterVariant: FilterVariant
): FilterOperator {
  const operators = getFilterOperators(filterVariant)
  return operators[0]?.value ?? (filterVariant === 'text' ? 'iLike' : 'eq')
}

export function getFilterOperators(
  filterVariant: FilterVariant
): { label: string; value: FilterOperator }[] {
  const operatorMap: Record<
    FilterVariant,
    { label: string; value: FilterOperator }[]
  > = {
    boolean: [
      { label: 'equals', value: 'eq' },
      { label: 'does not equal', value: 'ne' }
    ],
    date: [
      { label: 'equals', value: 'eq' },
      { label: 'does not equal', value: 'ne' },
      { label: 'before', value: 'lt' },
      { label: 'before or on', value: 'lte' },
      { label: 'after', value: 'gt' },
      { label: 'after or on', value: 'gte' },
      { label: 'is between', value: 'isBetween' },
      { label: 'is relative to today', value: 'isRelativeToToday' }
    ],
    multiSelect: [
      { label: 'is in', value: 'inArray' },
      { label: 'is not in', value: 'notInArray' }
    ],
    number: [
      { label: 'equals', value: 'eq' },
      { label: 'does not equal', value: 'ne' },
      { label: 'less than', value: 'lt' },
      { label: 'less than or equal', value: 'lte' },
      { label: 'greater than', value: 'gt' },
      { label: 'greater than or equal', value: 'gte' },
      { label: 'is between', value: 'isBetween' }
    ],
    select: [
      { label: 'equals', value: 'eq' },
      { label: 'does not equal', value: 'ne' },
      { label: 'is in', value: 'inArray' },
      { label: 'is not in', value: 'notInArray' }
    ],
    text: [
      { label: 'contains', value: 'iLike' },
      { label: 'does not contain', value: 'notILike' },
      { label: 'equals', value: 'eq' },
      { label: 'does not equal', value: 'ne' },
      { label: 'is empty', value: 'isEmpty' },
      { label: 'is not empty', value: 'isNotEmpty' }
    ]
  }

  return operatorMap[filterVariant] || operatorMap.text
}
