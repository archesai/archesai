import { z } from 'zod'

export const RunpodConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      token: z.ZodString
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({
    enabled: z.literal(false)
  }),
  z.object({
    enabled: z.literal(true),
    token: z.string().min(1)
  })
])
export type RunpodConfig = z.infer<typeof RunpodConfigSchema>

export const ServerConfigSchema: z.ZodObject<{
  cors: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    origins: z.ZodDefault<z.ZodString>
  }>
  docs: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    export: z.ZodDefault<z.ZodBoolean>
  }>
  host: z.ZodDefault<z.ZodString>
  port: z.ZodDefault<z.ZodNumber>
}> = z.object({
  cors: z.object({
    enabled: z.boolean().default(true),
    origins: z.string().default('localhost')
  }),
  docs: z.object({
    enabled: z.boolean().default(false),
    export: z.boolean().default(false)
  }),
  host: z.string().default('localhost'),
  port: z.number().default(3001)
})
export type ServerConfig = z.infer<typeof ServerConfigSchema>

export const ConfigConfigSchema: z.ZodObject<{
  validate: z.ZodDefault<z.ZodBoolean>
}> = z.object({
  validate: z.boolean().default(true)
})
export type ConfigConfig = z.infer<typeof ConfigConfigSchema>

export const PlatformConfigSchema: z.ZodObject<{
  enabled: z.ZodDefault<z.ZodBoolean>
  host: z.ZodDefault<z.ZodString>
}> = z.object({
  enabled: z.boolean().default(false),
  host: z.string().default('localhost')
})
export type PlatformConfig = z.infer<typeof PlatformConfigSchema>

export const TlsConfigSchema: z.ZodObject<{
  enabled: z.ZodDefault<z.ZodBoolean>
}> = z.object({
  enabled: z.boolean().default(false)
})
export type TlsConfig = z.infer<typeof TlsConfigSchema>

export const DatabaseConfigSchema: z.ZodObject<{
  type: z.ZodDefault<
    z.ZodEnum<{
      'in-memory': 'in-memory'
      postgres: 'postgres'
      sqlite: 'sqlite'
    }>
  >
  url: z.ZodDefault<z.ZodString>
}> = z.object({
  type: z.enum(['postgres', 'sqlite', 'in-memory']).default('postgres'),
  url: z.string().default('postgres://localhost:5432/arches')
})
export type DatabaseConfig = z.infer<typeof DatabaseConfigSchema>

export const EmailConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      password: z.ZodString
      service: z.ZodString
      user: z.ZodString
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({
    enabled: z.literal(false)
  }),
  z.object({
    enabled: z.literal(true),
    password: z.string(),
    service: z.string(),
    user: z.string()
  })
])
export type EmailConfig = z.infer<typeof EmailConfigSchema>

export const EmbeddingConfigSchema: z.ZodObject<{
  type: z.ZodDefault<
    z.ZodEnum<{
      ollama: 'ollama'
      openai: 'openai'
    }>
  >
}> = z.object({
  type: z.enum(['openai', 'ollama']).default('ollama')
})
export type EmbeddingConfig = z.infer<typeof EmbeddingConfigSchema>

export const SpeechConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      token: z.ZodString
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({ enabled: z.literal(false) }),
  z.object({
    enabled: z.literal(true),
    token: z.string()
  })
])
export type SpeechConfig = z.infer<typeof SpeechConfigSchema>

export const JwtConfigSchema: z.ZodObject<{
  expiration: z.ZodDefault<z.ZodString>
  secret: z.ZodDefault<z.ZodString>
}> = z.object({
  expiration: z.string().default((60 * 60 * 24).toString()),
  secret: z.string().default('secret-scary-stuff')
})
export type JwtConfig = z.infer<typeof JwtConfigSchema>

export const BillingConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      stripe: z.ZodObject<{
        token: z.ZodString
        whsec: z.ZodString
      }>
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({ enabled: z.literal(false) }),
  z.object({
    enabled: z.literal(true),
    stripe: z.object({ token: z.string(), whsec: z.string() })
  })
])
export type BillingConfig = z.infer<typeof BillingConfigSchema>

export const LlmConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      endpoint: z.ZodDefault<z.ZodString>
      token: z.ZodOptional<z.ZodString>
      type: z.ZodLiteral<'ollama'>
    }>,
    z.ZodObject<{
      endpoint: z.ZodOptional<z.ZodString>
      token: z.ZodString
      type: z.ZodLiteral<'openai'>
    }>
  ]
> = z.discriminatedUnion('type', [
  z.object({
    endpoint: z.string().default('http://localhost:11434'),
    token: z.string().optional(),
    type: z.literal('ollama')
  }),
  z.object({
    endpoint: z.string().optional(),
    token: z.string(),
    type: z.literal('openai')
  })
])
export type LlmConfig = z.infer<typeof LlmConfigSchema>

export const StorageConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      type: z.ZodLiteral<'local'>
    }>,
    z.ZodObject<{
      type: z.ZodLiteral<'google-cloud'>
    }>,
    z.ZodObject<{
      accesskey: z.ZodString
      bucket: z.ZodString
      endpoint: z.ZodString
      secretkey: z.ZodString
      type: z.ZodLiteral<'minio'>
    }>
  ]
> = z.discriminatedUnion('type', [
  z.object({ type: z.literal('local') }),
  z.object({ type: z.literal('google-cloud') }),
  z.object({
    accesskey: z.string(),
    bucket: z.string(),
    endpoint: z.string(),
    secretkey: z.string(),
    type: z.literal('minio')
  })
])
export type StorageConfig = z.infer<typeof StorageConfigSchema>

export const RedisConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      auth: z.ZodOptional<z.ZodString>
      ca: z.ZodOptional<z.ZodString>
      enabled: z.ZodLiteral<true>
      host: z.ZodString
      port: z.ZodNumber
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({ enabled: z.literal(false) }),
  z.object({
    auth: z.string().optional(),
    ca: z.string().optional(),
    enabled: z.literal(true),
    host: z.string(),
    port: z.number()
  })
])
export type RedisConfig = z.infer<typeof RedisConfigSchema>

export const SessionConfigSchema: z.ZodObject<{
  enabled: z.ZodDefault<z.ZodBoolean>
  secret: z.ZodDefault<z.ZodString>
}> = z.object({
  enabled: z.boolean().default(true),
  secret: z.string().default('session-scary-stuff')
})
export type SessionConfig = z.infer<typeof SessionConfigSchema>

export const MonitoringConfigSchema: z.ZodObject<{
  enabled: z.ZodDefault<z.ZodBoolean>
  loki: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    host: z.ZodOptional<z.ZodString>
  }>
}> = z.object({
  enabled: z.boolean().default(false),
  loki: z.object({
    enabled: z.boolean().default(false),
    host: z.string().optional()
  })
})
export type MonitoringConfig = z.infer<typeof MonitoringConfigSchema>

export const LoggingConfigSchema: z.ZodObject<{
  gcpfix: z.ZodDefault<z.ZodBoolean>
  level: z.ZodDefault<
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
  pretty: z.ZodDefault<z.ZodBoolean>
}> = z.object({
  gcpfix: z.boolean().default(false),
  level: z
    .enum(['fatal', 'error', 'warn', 'info', 'debug', 'trace', 'silent'])
    .default('info'),
  pretty: z.boolean().default(false)
})
export type LoggingConfig = z.infer<typeof LoggingConfigSchema>

export const ScraperConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      endpoint: z.ZodString
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({ enabled: z.literal(false) }),
  z.object({
    enabled: z.literal(true),
    endpoint: z.string()
  })
])
export type ScraperConfig = z.infer<typeof ScraperConfigSchema>

export const UnstructuredConfigSchema: z.ZodDiscriminatedUnion<
  [
    z.ZodObject<{
      enabled: z.ZodLiteral<false>
    }>,
    z.ZodObject<{
      enabled: z.ZodLiteral<true>
      endpoint: z.ZodString
    }>
  ]
> = z.discriminatedUnion('enabled', [
  z.object({ enabled: z.literal(false) }),
  z.object({
    enabled: z.literal(true),
    endpoint: z.string()
  })
])
export type UnstructuredConfig = z.infer<typeof UnstructuredConfigSchema>

export const AuthConfigSchema: z.ZodObject<{
  firebase: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        clientEmail: z.ZodString
        enabled: z.ZodLiteral<true>
        privateKey: z.ZodString
        projectId: z.ZodString
      }>
    ]
  >
  local: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
  }>
  twitter: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        callbackURL: z.ZodString
        consumerKey: z.ZodString
        consumerSecret: z.ZodString
        enabled: z.ZodLiteral<true>
      }>
    ]
  >
}> = z.object({
  firebase: z.discriminatedUnion('enabled', [
    z.object({ enabled: z.literal(false) }),
    z.object({
      clientEmail: z.string(),
      enabled: z.literal(true),
      privateKey: z.string(),
      projectId: z.string()
    })
  ]),
  local: z.object({
    enabled: z.boolean().default(false)
  }),
  twitter: z.discriminatedUnion('enabled', [
    z.object({ enabled: z.literal(false) }),
    z.object({
      callbackURL: z.string(),
      consumerKey: z.string(),
      consumerSecret: z.string(),
      enabled: z.literal(true)
    })
  ])
})
export type AuthConfig = z.infer<typeof AuthConfigSchema>

export const ArchesConfigSchema: z.ZodObject<{
  auth: z.ZodObject<{
    firebase: z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          enabled: z.ZodLiteral<false>
        }>,
        z.ZodObject<{
          clientEmail: z.ZodString
          enabled: z.ZodLiteral<true>
          privateKey: z.ZodString
          projectId: z.ZodString
        }>
      ]
    >
    local: z.ZodObject<{
      enabled: z.ZodDefault<z.ZodBoolean>
    }>
    twitter: z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          enabled: z.ZodLiteral<false>
        }>,
        z.ZodObject<{
          callbackURL: z.ZodString
          consumerKey: z.ZodString
          consumerSecret: z.ZodString
          enabled: z.ZodLiteral<true>
        }>
      ]
    >
  }>
  billing: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        stripe: z.ZodObject<{
          token: z.ZodString
          whsec: z.ZodString
        }>
      }>
    ]
  >
  config: z.ZodObject<{
    validate: z.ZodDefault<z.ZodBoolean>
  }>
  database: z.ZodObject<{
    type: z.ZodDefault<
      z.ZodEnum<{
        'in-memory': 'in-memory'
        postgres: 'postgres'
        sqlite: 'sqlite'
      }>
    >
    url: z.ZodDefault<z.ZodString>
  }>
  email: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        password: z.ZodString
        service: z.ZodString
        user: z.ZodString
      }>
    ]
  >
  embedding: z.ZodObject<{
    type: z.ZodDefault<
      z.ZodEnum<{
        ollama: 'ollama'
        openai: 'openai'
      }>
    >
  }>
  jwt: z.ZodObject<{
    expiration: z.ZodDefault<z.ZodString>
    secret: z.ZodDefault<z.ZodString>
  }>
  llm: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        endpoint: z.ZodDefault<z.ZodString>
        token: z.ZodOptional<z.ZodString>
        type: z.ZodLiteral<'ollama'>
      }>,
      z.ZodObject<{
        endpoint: z.ZodOptional<z.ZodString>
        token: z.ZodString
        type: z.ZodLiteral<'openai'>
      }>
    ]
  >
  logging: z.ZodObject<{
    gcpfix: z.ZodDefault<z.ZodBoolean>
    level: z.ZodDefault<
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
    pretty: z.ZodDefault<z.ZodBoolean>
  }>
  monitoring: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    loki: z.ZodObject<{
      enabled: z.ZodDefault<z.ZodBoolean>
      host: z.ZodOptional<z.ZodString>
    }>
  }>
  platform: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    host: z.ZodDefault<z.ZodString>
  }>
  redis: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        auth: z.ZodOptional<z.ZodString>
        ca: z.ZodOptional<z.ZodString>
        enabled: z.ZodLiteral<true>
        host: z.ZodString
        port: z.ZodNumber
      }>
    ]
  >
  runpod: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        token: z.ZodString
      }>
    ]
  >
  scraper: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        endpoint: z.ZodString
      }>
    ]
  >
  server: z.ZodObject<{
    cors: z.ZodObject<{
      enabled: z.ZodDefault<z.ZodBoolean>
      origins: z.ZodDefault<z.ZodString>
    }>
    docs: z.ZodObject<{
      enabled: z.ZodDefault<z.ZodBoolean>
      export: z.ZodDefault<z.ZodBoolean>
    }>
    host: z.ZodDefault<z.ZodString>
    port: z.ZodDefault<z.ZodNumber>
  }>
  session: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
    secret: z.ZodDefault<z.ZodString>
  }>
  speech: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        token: z.ZodString
      }>
    ]
  >
  storage: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        type: z.ZodLiteral<'local'>
      }>,
      z.ZodObject<{
        type: z.ZodLiteral<'google-cloud'>
      }>,
      z.ZodObject<{
        accesskey: z.ZodString
        bucket: z.ZodString
        endpoint: z.ZodString
        secretkey: z.ZodString
        type: z.ZodLiteral<'minio'>
      }>
    ]
  >
  tls: z.ZodObject<{
    enabled: z.ZodDefault<z.ZodBoolean>
  }>
  unstructured: z.ZodDiscriminatedUnion<
    [
      z.ZodObject<{
        enabled: z.ZodLiteral<false>
      }>,
      z.ZodObject<{
        enabled: z.ZodLiteral<true>
        endpoint: z.ZodString
      }>
    ]
  >
}> = z
  .object({
    auth: AuthConfigSchema,
    billing: BillingConfigSchema,
    config: ConfigConfigSchema,
    database: DatabaseConfigSchema,
    email: EmailConfigSchema,
    embedding: EmbeddingConfigSchema,
    jwt: JwtConfigSchema,
    llm: LlmConfigSchema,
    logging: LoggingConfigSchema,
    monitoring: MonitoringConfigSchema,
    platform: PlatformConfigSchema,
    redis: RedisConfigSchema,
    runpod: RunpodConfigSchema,
    scraper: ScraperConfigSchema,
    server: ServerConfigSchema,
    session: SessionConfigSchema,
    speech: SpeechConfigSchema,
    storage: StorageConfigSchema,
    tls: TlsConfigSchema,
    unstructured: UnstructuredConfigSchema
  })
  .describe('Arches AI configuration schema')
export type ArchesConfig = z.infer<typeof ArchesConfigSchema>
