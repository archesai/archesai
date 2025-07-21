import { z } from 'zod'

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
