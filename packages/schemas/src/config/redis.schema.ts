import { z } from 'zod'

const BaseRedisConfigSchema = z.object({
  auth: z
    .string()
    .optional()
    .default('password')
    .describe('Redis authentication password (optional)'),
  ca: z
    .string()
    .optional()
    .describe('Certificate Authority for TLS connections (optional)'),
  host: z
    .string()
    .optional()
    .default('localhost')
    .describe('Redis server hostname or IP address'),
  port: z.number().optional().default(6379).describe('Redis server port number')
})

export const RedisConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          auth: z.ZodDefault<z.ZodOptional<z.ZodString>>
          ca: z.ZodOptional<z.ZodString>
          host: z.ZodDefault<z.ZodOptional<z.ZodString>>
          mode: z.ZodLiteral<'enabled'>
          port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
        }>,
        z.ZodObject<{
          auth: z.ZodDefault<z.ZodOptional<z.ZodString>>
          ca: z.ZodOptional<z.ZodString>
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
          persistence: z.ZodObject<{
            enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
            size: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
          port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
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
    BaseRedisConfigSchema.extend({
      mode: z.literal('enabled')
    }),
    BaseRedisConfigSchema.extend({
      image: z.object({
        pullPolicy: z
          .enum(['Always', 'IfNotPresent', 'Never'])
          .optional()
          .default('IfNotPresent'),
        repository: z.string().optional().default('redis'),
        tag: z.string().optional().default('7-alpine')
      }),
      mode: z.literal('managed'),
      persistence: z.object({
        enabled: z
          .boolean()
          .optional()
          .default(true)
          .describe('Enable persistent storage for Redis data'),
        size: z
          .string()
          .optional()
          .default('1Gi')
          .describe('Size of persistent volume for Redis storage')
      }),
      resources: z.object({
        limits: z.object({
          cpu: z.string().optional().default('200m'),
          memory: z.string().optional().default('256Mi')
        }),
        requests: z.object({
          cpu: z.string().optional().default('100m'),
          memory: z.string().optional().default('128Mi')
        })
      })
    })
  ])
  .optional()
  .default({
    auth: 'password',
    host: 'localhost',
    image: {
      pullPolicy: 'IfNotPresent',
      repository: 'redis',
      tag: '7-alpine'
    },
    mode: 'managed',
    persistence: {
      enabled: true,
      size: '1Gi'
    },
    port: 6379,
    resources: {
      limits: { cpu: '200m', memory: '256Mi' },
      requests: { cpu: '100m', memory: '128Mi' }
    }
  })

export type RedisConfig = z.infer<typeof RedisConfigSchema>
