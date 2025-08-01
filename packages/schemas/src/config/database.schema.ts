import { z } from 'zod'

const BaseDatabaseConfigSchema = z.object({
  url: z
    .string()
    .optional()
    .default(
      'postgresql://admin:password@localhost:5432/archesai-db?schema=public'
    )
    .describe('Database connection URL/string')
})

export const DatabaseConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          mode: z.ZodLiteral<'enabled'>
          url: z.ZodDefault<z.ZodOptional<z.ZodString>>
        }>,
        z.ZodObject<{
          auth: z.ZodDefault<
            z.ZodOptional<
              z.ZodObject<{
                database: z.ZodDefault<z.ZodOptional<z.ZodString>>
                password: z.ZodDefault<z.ZodOptional<z.ZodString>>
              }>
            >
          >
          image: z.ZodDefault<
            z.ZodOptional<
              z.ZodObject<{
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
            >
          >
          mode: z.ZodLiteral<'managed'>
          persistence: z.ZodObject<{
            enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
            size: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
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
    BaseDatabaseConfigSchema.extend({
      mode: z.literal('enabled')
    }),
    BaseDatabaseConfigSchema.extend({
      auth: z
        .object({
          database: z
            .string()
            .optional()
            .default('archesai-db')
            .describe('Database name to create and use'),
          password: z
            .string()
            .optional()
            .default('password')
            .describe('Database user password (change for production!)')
        })
        .optional()
        .default({
          database: 'archesai-db',
          password: 'password'
        }),
      image: z
        .object({
          pullPolicy: z
            .enum(['Always', 'IfNotPresent', 'Never'])
            .optional()
            .default('IfNotPresent')
            .describe('Kubernetes image pull policy'),
          repository: z
            .string()
            .optional()
            .default('pgvector/pgvector')
            .describe('PostgreSQL with pgvector extension docker image'),
          tag: z
            .string()
            .optional()
            .default('pg16')
            .describe('PostgreSQL version tag')
        })
        .optional()
        .default({
          pullPolicy: 'IfNotPresent',
          repository: 'pgvector/pgvector',
          tag: 'pg16'
        }),
      mode: z.literal('managed'),
      persistence: z.object({
        enabled: z
          .boolean()
          .optional()
          .default(true)
          .describe('Enable persistent storage for database data'),
        size: z
          .string()
          .optional()
          .default('10Gi')
          .describe('Size of persistent volume for database storage')
      }),
      resources: z.object({
        limits: z.object({
          cpu: z
            .string()
            .optional()
            .default('500m')
            .describe('Maximum CPU allocation for database'),
          memory: z
            .string()
            .optional()
            .default('1Gi')
            .describe('Maximum memory allocation for database')
        }),
        requests: z.object({
          cpu: z
            .string()
            .optional()
            .default('250m')
            .describe('Requested CPU allocation for database'),
          memory: z
            .string()
            .optional()
            .default('512Mi')
            .describe('Requested memory allocation for database')
        })
      })
    })
  ])
  .optional()
  .default({
    auth: {
      database: 'archesai-db',
      password: 'password'
    },
    image: {
      pullPolicy: 'IfNotPresent',
      repository: 'pgvector/pgvector',
      tag: 'pg16'
    },
    mode: 'managed',
    persistence: {
      enabled: true,
      size: '10Gi'
    },
    resources: {
      limits: { cpu: '500m', memory: '1Gi' },
      requests: { cpu: '250m', memory: '512Mi' }
    },
    url: 'postgresql://admin:password@localhost:5432/archesai-db?schema=public'
  })
  .describe(
    'Database configuration for PostgreSQL with optional pgvector support. Includes managed mode with persistence and resource limits.'
  )

export type DatabaseConfig = z.infer<typeof DatabaseConfigSchema>
