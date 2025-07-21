import { z } from 'zod'

export const RecordSchema: z.ZodRecord<z.ZodString, z.ZodUnknown> = z.record(
  z.string(),
  z.unknown()
)
