import { z } from 'zod'

import { FilterValueSchema } from '#query/filter-value.schema'
import { OperatorSchema } from '#query/operator.schema'

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

export type FilterConditionType = z.infer<typeof FilterConditionSchema>
