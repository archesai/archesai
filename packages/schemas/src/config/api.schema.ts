import { z } from 'zod'

const BaseApiConfigSchema = z.object({
  cors: z
    .object({
      origins: z
        .string()
        .optional()
        .default('https://platform.archesai.dev')
        .describe(
          'A comma-separated list of allowed origins for CORS requests. Use "*" to allow all'
        )
    })
    .optional()
    .default({ origins: 'https://platform.archesai.dev' })
    .describe(
      'CORS configuration for the API server. This allows you to specify which origins are allowed to make requests to the API.'
    ),
  docs: z
    .boolean()
    .optional()
    .default(true)
    .describe('Enable or disable API documentation'),
  email: z
    .discriminatedUnion('mode', [
      z
        .object({
          mode: z.literal('disabled')
        })
        .describe('Email configuration is disabled'),
      z
        .object({
          mode: z.literal('enabled').describe('Email configuration is enabled'),
          password: z
            .string()
            .describe(
              'Password for the email service. This is required when email configuration is enabled.'
            ),
          service: z
            .string()
            .describe(
              'Email service provider (e.g., "gmail", "sendgrid", etc.). This is required when email configuration is enabled.'
            ),
          user: z
            .string()
            .describe(
              'Username for the email service. This is required when email configuration is enabled.'
            )
        })
        .describe(
          'Email configuration for sending emails. This includes the service, user, and password for the email service.'
        )
    ])
    .optional()
    .default({ mode: 'disabled' }),
  host: z
    .string()
    .optional()
    .default('0.0.0.0')
    .describe('The host address on which the API server will listen'),
  port: z
    .number()
    .optional()
    .default(3001)
    .describe('The port on which the API server will listen'),
  validate: z
    .boolean()
    .optional()
    .default(true)
    .describe(
      'Enable or disable request validation. When enabled, the API will validate incoming requests against the defined schemas.'
    )
})

export const ApiConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodObject<{
      cors: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            origins: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
        >
      >
      docs: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
      email: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                mode: z.ZodLiteral<'enabled'>
                password: z.ZodString
                service: z.ZodString
                user: z.ZodString
              }>
            ]
          >
        >
      >
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
      validate: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
    }>
  >
> = BaseApiConfigSchema.extend({
  image: z.object({
    pullPolicy: z
      .enum(['Always', 'IfNotPresent', 'Never'])
      .optional()
      .default('IfNotPresent'),
    repository: z.string().optional().default('archesai/api'),
    tag: z.string().optional().default('latest')
  }),
  resources: z.object({
    limits: z.object({
      cpu: z.string().optional().default('1000m'),
      memory: z.string().optional().default('1Gi')
    }),
    requests: z.object({
      cpu: z.string().optional().default('500m'),
      memory: z.string().optional().default('512Mi')
    })
  })
})
  .optional()
  .default({
    cors: { origins: 'https://platform.archesai.dev' },
    docs: true,
    email: { mode: 'disabled' },
    host: '0.0.0.0',
    image: {
      pullPolicy: 'IfNotPresent',
      repository: 'archesai/api',
      tag: 'latest'
    },
    port: 3001,
    resources: {
      limits: { cpu: '1000m', memory: '1Gi' },
      requests: { cpu: '500m', memory: '512Mi' }
    },
    validate: true
  })
  .describe(
    'Configuration schema for the API server. This includes settings for CORS, documentation, email, host, port, and request validation.'
  )

export type ApiConfig = z.infer<typeof ApiConfigSchema>
