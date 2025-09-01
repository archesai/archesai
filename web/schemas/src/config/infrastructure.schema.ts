import { z } from 'zod'

export const InfrastructureConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodObject<{
      development: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            api: z.ZodDefault<
              z.ZodOptional<
                z.ZodObject<{
                  enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                  port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
                }>
              >
            >
            hostIP: z.ZodDefault<z.ZodOptional<z.ZodString>>
            loki: z.ZodDefault<
              z.ZodOptional<
                z.ZodObject<{
                  enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                  port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
                }>
              >
            >
            platform: z.ZodDefault<
              z.ZodOptional<
                z.ZodObject<{
                  enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                  port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
                }>
              >
            >
            postgres: z.ZodDefault<
              z.ZodOptional<
                z.ZodObject<{
                  enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                  port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
                }>
              >
            >
            redis: z.ZodDefault<
              z.ZodOptional<
                z.ZodObject<{
                  enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                  port: z.ZodDefault<z.ZodOptional<z.ZodNumber>>
                }>
              >
            >
          }>
        >
      >
      images: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            imagePullSecrets: z.ZodDefault<
              z.ZodOptional<z.ZodArray<z.ZodString>>
            >
            imageRegistry: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
        >
      >
      migrations: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
          }>
        >
      >
      namespace: z.ZodDefault<z.ZodOptional<z.ZodString>>
      serviceAccount: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            create: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
            name: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
        >
      >
    }>
  >
> = z
  .object({
    development: z
      .object({
        api: z
          .object({
            enabled: z
              .boolean()
              .optional()
              .default(false)
              .describe(
                'Enable development port forwarding for the API service'
              ),
            port: z
              .number()
              .optional()
              .default(3001)
              .describe(
                'Local port to forward API service to during development'
              )
          })
          .optional()
          .default({
            enabled: false,
            port: 3001
          })
          .describe(
            'Development port forwarding configuration for the API service'
          ),
        hostIP: z
          .string()
          .optional()
          .default('172.18.0.1')
          .describe(
            'Host IP address for development port forwarding (typically Docker bridge IP)'
          ),
        loki: z
          .object({
            enabled: z
              .boolean()
              .optional()
              .default(false)
              .describe('Enable development port forwarding for Loki service'),
            port: z
              .number()
              .optional()
              .default(30056)
              .describe(
                'Local port to forward Loki service to during development'
              )
          })
          .optional()
          .default({
            enabled: false,
            port: 30056
          })
          .describe(
            'Development port forwarding configuration for the Loki service'
          ),
        platform: z
          .object({
            enabled: z
              .boolean()
              .optional()
              .default(false)
              .describe(
                'Enable development port forwarding for the platform/frontend service'
              ),
            port: z
              .number()
              .optional()
              .default(3000)
              .describe(
                'Local port to forward platform service to during development'
              )
          })
          .optional()
          .default({
            enabled: false,
            port: 3000
          })
          .describe(
            'Development port forwarding configuration for the platform service'
          ),
        postgres: z
          .object({
            enabled: z
              .boolean()
              .optional()
              .default(false)
              .describe(
                'Enable development port forwarding for PostgreSQL database'
              ),
            port: z
              .number()
              .optional()
              .default(30054)
              .describe(
                'Local port to forward PostgreSQL to during development'
              )
          })
          .optional()
          .default({
            enabled: false,
            port: 30054
          })
          .describe(
            'Development port forwarding configuration for PostgreSQL database'
          ),
        redis: z
          .object({
            enabled: z
              .boolean()
              .optional()
              .default(false)
              .describe('Enable development port forwarding for Redis cache'),
            port: z
              .number()
              .optional()
              .default(30055)
              .describe('Local port to forward Redis to during development')
          })
          .optional()
          .default({
            enabled: false,
            port: 30055
          })
          .describe('Development port forwarding configuration for Redis cache')
      })
      .optional()
      .default({
        api: { enabled: false, port: 3001 },
        hostIP: '172.18.0.1',
        loki: { enabled: false, port: 30056 },
        platform: { enabled: false, port: 3000 },
        postgres: { enabled: false, port: 30054 },
        redis: { enabled: false, port: 30055 }
      })
      .describe(
        'Development environment configuration for local port forwarding and debugging'
      ),
    images: z
      .object({
        imagePullSecrets: z
          .array(z.string())
          .optional()
          .default([])
          .describe(
            'List of Kubernetes secrets for pulling private container images'
          ),

        imageRegistry: z
          .string()
          .optional()
          .default('')
          .describe(
            'Custom container registry URL (leave empty for Docker Hub)'
          )
      })
      .optional()
      .default({
        imagePullSecrets: [],
        imageRegistry: ''
      })
      .describe('Container image configuration for Kubernetes deployments'),

    migrations: z
      .object({
        enabled: z
          .boolean()
          .optional()
          .default(false)
          .describe('Enable automatic database migrations on deployment')
      })
      .optional()
      .default({ enabled: false })
      .describe('Database migration configuration for schema updates'),

    namespace: z
      .string()
      .optional()
      .default('arches-system')
      .describe('Kubernetes namespace where all resources will be deployed'),

    serviceAccount: z
      .object({
        create: z
          .boolean()
          .optional()
          .default(true)
          .describe(
            'Create a dedicated Kubernetes service account for the application'
          ),
        name: z
          .string()
          .optional()
          .default('')
          .describe('Custom service account name (auto-generated if empty)')
      })
      .optional()
      .default({ create: true, name: '' })
      .describe(
        'Kubernetes service account configuration for pod security and RBAC'
      )
  })
  .optional()
  .default({
    development: {
      api: { enabled: false, port: 3001 },
      hostIP: '172.18.0.1',
      loki: { enabled: false, port: 30056 },
      platform: { enabled: false, port: 3000 },
      postgres: { enabled: false, port: 30054 },
      redis: { enabled: false, port: 30055 }
    },
    images: {
      imagePullSecrets: [],
      imageRegistry: ''
    },
    migrations: { enabled: false },
    namespace: 'arches-system',
    serviceAccount: { create: true, name: '' }
  })
  .describe(
    'Infrastructure configuration for Kubernetes deployments, including development settings, image management, migrations, and service accounts.'
  )

export type InfrastructureConfig = z.infer<typeof InfrastructureConfigSchema>
