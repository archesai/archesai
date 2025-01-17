import { z } from 'zod'

export type ArchesConfig = z.infer<typeof archesConfigSchema>

const stringToBoolean = z.preprocess((val) => {
  if (typeof val === 'string') {
    if (val.toLowerCase() === 'true') return true
    if (val.toLowerCase() === 'false') return false
  }
  return val
}, z.boolean())

export const archesConfigSchema = z.object({
  server: z.object({
    host: z.string(),
    cors: z.object({
      origins: z.string()
    }),
    docs: z.object({
      enabled: stringToBoolean
    })
  }),

  config: z.object({
    validate: stringToBoolean
  }),

  frontend: z.object({
    host: z.string()
  }),

  tls: z.object({
    enabled: stringToBoolean
  }),

  database: z.object({
    url: z.string()
  }),

  email: z
    .object({
      enabled: stringToBoolean,
      user: z.string().optional(),
      password: z.string().optional(),
      service: z.string().optional()
    })
    .refine(
      (data) => {
        if (data.enabled) {
          return !!data.user && !!data.password && !!data.service
        }
        return true
      },
      {
        message:
          'Email user, password, and service are required when email.enabled is true.',
        path: ['email.user', 'email.password', 'email.service'] // Specify paths for better error reporting
      }
    ),

  embedding: z.object({
    type: z.enum(['openai', 'ollama'])
  }),

  jwt: z.object({
    expiration: z.string(),
    secret: z.string()
  }),

  billing: z
    .object({
      enabled: stringToBoolean,
      stripe: z
        .object({
          whsec: z.string(),
          token: z.string()
        })
        .optional()
    })
    .refine(
      (data) => {
        if (data.enabled) {
          return !!data.stripe?.token && !!data.stripe?.whsec
        }
        return true
      },
      {
        message:
          'Stripe private API key and webhook secret are required when billing.enabled is true.',
        path: ['billing.stripe.whsec', 'billing.stripe.token']
      }
    ),

  llm: z
    .object({
      type: z.enum(['openai', 'ollama']),
      endpoint: z.string().optional(),
      token: z.string().optional()
    })
    .refine(
      (data) => {
        if (data.type === 'ollama') {
          return !!data.endpoint
        }
        if (data.type === 'openai') {
          return !!data.token
        }
        return true
      },
      {
        message:
          'llm.endpoint is required for ollama, and llm.token is required for openai.',
        path: ['llm.endpoint', 'llm.token']
      }
    ),

  storage: z
    .object({
      type: z.enum(['google-cloud', 'local', 'minio']),
      endpoint: z.string().optional(),
      accesskey: z.string().optional(),
      secretkey: z.string().optional(),
      bucket: z.string().optional()
    })
    .refine(
      (data) => {
        if (data.type === 'minio') {
          return (
            !!data.endpoint &&
            !!data.accesskey &&
            !!data.secretkey &&
            !!data.bucket
          )
        }
        return true
      },
      {
        message:
          'Minio host, access key, secret key, and bucket are required when storage.type is minio.',
        path: [
          'storage.endpoint',
          'storage.accesskey',
          'storage.secretkey',
          'storage.bucket'
        ]
      }
    ),

  redis: z.object({
    auth: z.string(),
    ca: z.string().optional(),
    host: z.string(),
    port: z.coerce.number()
  }),

  session: z.object({
    secret: z.string()
  }),

  monitoring: z.object({
    enabled: stringToBoolean,
    loki: z
      .object({
        enabled: stringToBoolean,
        host: z.string()
      })
      .optional()
  }),

  logging: z.object({
    level: z.enum([
      'fatal',
      'error',
      'warn',
      'info',
      'debug',
      'trace',
      'silent'
    ]),
    pretty: stringToBoolean,
    gcpfix: stringToBoolean.optional().default(false)
  }),

  scraper: z
    .object({
      enabled: stringToBoolean,
      endpoint: z.string().optional()
    })
    .refine(
      (data) => {
        if (data.enabled) {
          return !!data.endpoint
        }
        return true
      },
      {
        message: 'scraper.endpoint is required when scraper.enabled is true.',
        path: ['scraper.enabled']
      }
    ),

  unstructured: z
    .object({
      enabled: stringToBoolean,
      endpoint: z.string().optional()
    })
    .refine(
      (data) => {
        if (data.enabled) {
          return !!data.endpoint
        }
        return true
      },
      {
        message:
          'unstructured.endpoint is required when unstructured.enabled is true.',
        path: ['unstructured.endpoint']
      }
    )
})
