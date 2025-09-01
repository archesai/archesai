import { z } from 'zod'

export const FilterOperationEnum = [
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

export type FilterOperation = (typeof FilterOperationEnum)[number]

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
}> = z.enum(FilterOperationEnum).meta({
  description: 'Supported filter operations',
  id: 'Operator'
})
