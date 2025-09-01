import { z } from 'zod'

const BasePlatformConfigSchema = z.object({
  host: z
    .string()
    .optional()
    .default('localhost')
    .describe('Host address where the platform service will be accessible')
})

export const PlatformConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          host: z.ZodDefault<z.ZodOptional<z.ZodString>>
          mode: z.ZodLiteral<'enabled'>
        }>,
        z.ZodObject<{
          host: z.ZodDefault<z.ZodOptional<z.ZodString>>
          image: z.ZodObject<{
            pullPolicy: z.ZodDefault<
              z.ZodOptional<
                z.ZodEnum<{
                  Always: 'Always'
                  IfNotPresent: 'IfNotPresent'
                  Never: 'Never'
                }>
              >
            >
            repository: z.ZodDefault<z.ZodOptional<z.ZodString>>
            tag: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
          mode: z.ZodLiteral<'managed'>
          resources: z.ZodObject<{
            limits: z.ZodObject<{
              cpu: z.ZodDefault<z.ZodOptional<z.ZodString>>
              memory: z.ZodDefault<z.ZodOptional<z.ZodString>>
            }>
            requests: z.ZodObject<{
              cpu: z.ZodDefault<z.ZodOptional<z.ZodString>>
              memory: z.ZodDefault<z.ZodOptional<z.ZodString>>
            }>
          }>
        }>
      ]
    >
  >
> = z
  .discriminatedUnion('mode', [
    z.object({
      mode: z.literal('disabled')
    }),
    BasePlatformConfigSchema.extend({
      mode: z.literal('enabled')
    }),
    BasePlatformConfigSchema.extend({
      image: z.object({
        pullPolicy: z
          .enum(['Always', 'IfNotPresent', 'Never'])
          .optional()
          .default('IfNotPresent'),
        repository: z.string().optional().default('archesai/platform'),
        tag: z.string().optional().default('latest')
      }),
      mode: z.literal('managed'),
      resources: z.object({
        limits: z.object({
          cpu: z.string().optional().default('500m'),
          memory: z.string().optional().default('512Mi')
        }),
        requests: z.object({
          cpu: z.string().optional().default('250m'),
          memory: z.string().optional().default('256Mi')
        })
      })
    })
  ])
  .optional()
  .default({
    host: 'localhost',
    image: {
      pullPolicy: 'IfNotPresent',
      repository: 'archesai/platform',
      tag: 'latest'
    },
    mode: 'managed',
    resources: {
      limits: { cpu: '500m', memory: '512Mi' },
      requests: { cpu: '250m', memory: '256Mi' }
    }
  })
  .describe(
    'Platform configuration for the Arches AI platform service, including host address, image settings, and resource limits.'
  )

export type PlatformConfig = z.infer<typeof PlatformConfigSchema>
