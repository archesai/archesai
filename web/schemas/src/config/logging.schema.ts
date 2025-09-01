import { z } from 'zod'

export const LoggingConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodObject<{
      level: z.ZodDefault<
        z.ZodOptional<
          z.ZodEnum<{
            debug: 'debug'
            error: 'error'
            fatal: 'fatal'
            info: 'info'
            silent: 'silent'
            trace: 'trace'
            warn: 'warn'
          }>
        >
      >
      pretty: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
    }>
  >
> = z
  .object({
    level: z
      .enum(['fatal', 'error', 'warn', 'info', 'debug', 'trace', 'silent'])
      .optional()
      .default('info')
      .describe('Minimum log level to output (fatal=highest, silent=no logs)'),
    pretty: z
      .boolean()
      .optional()
      .default(false)
      .describe(
        'Enable pretty-printed logs for development (disable in production for structured logs)'
      )
  })
  .optional()
  .default({
    level: 'info',
    pretty: false
  })
  .describe(
    'Logging configuration for the application. This includes the log level and whether to pretty-print logs.'
  )

export type LoggingConfig = z.infer<typeof LoggingConfigSchema>
