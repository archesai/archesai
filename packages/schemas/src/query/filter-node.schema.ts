import { z } from 'zod'

import type { FilterConditionType } from '#query/filter-condition.schema'

import { FilterConditionSchema } from '#query/filter-condition.schema'

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
