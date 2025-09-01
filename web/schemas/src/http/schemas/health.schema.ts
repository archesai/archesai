import { z } from 'zod'

export const HealthCheckSchema: z.ZodObject<{
  services: z.ZodObject<{
    database: z.ZodString
    email: z.ZodString
    redis: z.ZodString
  }>
  timestamp: z.ZodString
  uptime: z.ZodNumber
}> = z.object({
  services: z.object({
    database: z.string(),
    email: z.string(),
    redis: z.string()
  }),
  timestamp: z.string(),
  uptime: z.number()
})
