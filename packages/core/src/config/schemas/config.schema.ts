import type {
  Static,
  TBoolean,
  TLiteral,
  TNumber,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const RunpodConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      token: TString
    }>
  ]
> = Type.Union([
  Type.Object({
    enabled: Type.Literal(false)
  }),
  Type.Object({
    enabled: Type.Literal(true),
    token: Type.String({ minLength: 1 })
  })
])
export type RunpodConfig = Static<typeof RunpodConfigSchema>

export const ServerConfigSchema: TObject<{
  cors: TObject<{
    enabled: TBoolean
    origins: TString
  }>
  docs: TObject<{
    enabled: TBoolean
    export: TBoolean
  }>
  host: TString
  port: TNumber
}> = Type.Object({
  cors: Type.Object({
    enabled: Type.Boolean({ default: true }),
    origins: Type.String({ default: 'localhost' })
  }),
  docs: Type.Object({
    enabled: Type.Boolean({ default: false }),
    export: Type.Boolean({ default: false })
  }),
  host: Type.String({ default: 'localhost' }),
  port: Type.Number({ default: 3001 })
})
export type ServerConfig = Static<typeof ServerConfigSchema>

export const ConfigConfigSchema: TObject<{
  validate: TBoolean
}> = Type.Object({
  validate: Type.Boolean({ default: true })
})
export type ConfigConfig = Static<typeof ConfigConfigSchema>

export const PlatformConfigSchema: TObject<{
  enabled: TBoolean
  host: TString
}> = Type.Object({
  enabled: Type.Boolean({ default: false }),
  host: Type.String({ default: 'localhost' })
})
export type PlatformConfig = Static<typeof PlatformConfigSchema>

export const TlsConfigSchema: TObject<{
  enabled: TBoolean
}> = Type.Object({
  enabled: Type.Boolean({ default: false })
})
export type TlsConfig = Static<typeof TlsConfigSchema>

export const DatabaseConfigSchema: TObject<{
  type: TUnion<
    [TLiteral<'postgres'>, TLiteral<'sqlite'>, TLiteral<'in-memory'>]
  >
  url: TString
}> = Type.Object({
  type: Type.Union(
    [
      Type.Literal('postgres'),
      Type.Literal('sqlite'),
      Type.Literal('in-memory')
    ],
    {
      default: 'postgres'
    }
  ),
  url: Type.String({ default: 'postgres://localhost:5432/arches' })
})
export type DatabaseConfig = Static<typeof DatabaseConfigSchema>

export const EmailConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      password: TString
      service: TString
      user: TString
    }>
  ]
> = Type.Union([
  Type.Object({
    enabled: Type.Literal(false)
  }),
  Type.Object({
    enabled: Type.Literal(true),
    password: Type.String(),
    service: Type.String(),
    user: Type.String()
  })
])
export type EmailConfig = Static<typeof EmailConfigSchema>

export const EmbeddingConfigSchema: TObject<{
  type: TUnion<[TLiteral<'openai'>, TLiteral<'ollama'>]>
}> = Type.Object({
  type: Type.Union([Type.Literal('openai'), Type.Literal('ollama')], {
    default: 'ollama'
  })
})
export type EmbeddingConfig = Static<typeof EmbeddingConfigSchema>

export const SpeechConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      token: TString
    }>
  ]
> = Type.Union([
  Type.Object({ enabled: Type.Literal(false) }),
  Type.Object({
    enabled: Type.Literal(true),
    token: Type.String()
  })
])
export type SpeechConfig = Static<typeof SpeechConfigSchema>

export const JwtConfigSchema: TObject<{
  expiration: TString
  secret: TString
}> = Type.Object({
  expiration: Type.String({ default: (60 * 60 * 24).toString() }),
  secret: Type.String({ default: 'secret-scary-stuff' })
})
export type JwtConfig = Static<typeof JwtConfigSchema>

export const BillingConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      stripe: TObject<{
        token: TString
        whsec: TString
      }>
    }>
  ]
> = Type.Union([
  Type.Object({ enabled: Type.Literal(false) }),
  Type.Object({
    enabled: Type.Literal(true),
    stripe: Type.Object({ token: Type.String(), whsec: Type.String() })
  })
])
export type BillingConfig = Static<typeof BillingConfigSchema>

export const LlmConfigSchema: TUnion<
  [
    TObject<{
      endpoint: TString
      token: TOptional<TString>
      type: TLiteral<'ollama'>
    }>,
    TObject<{
      endpoint: TOptional<TString>
      token: TString
      type: TLiteral<'openai'>
    }>
  ]
> = Type.Union([
  Type.Object({
    endpoint: Type.String({ default: 'http://localhost:11434' }),
    token: Type.Optional(Type.String()),
    type: Type.Literal('ollama')
  }),
  Type.Object({
    endpoint: Type.Optional(Type.String()),
    token: Type.String(),
    type: Type.Literal('openai')
  })
])
export type LlmConfig = Static<typeof LlmConfigSchema>

export const StorageConfigSchema: TUnion<
  [
    TObject<{
      type: TLiteral<'local'>
    }>,
    TObject<{
      type: TLiteral<'google-cloud'>
    }>,
    TObject<{
      accesskey: TString
      bucket: TString
      endpoint: TString
      secretkey: TString
      type: TLiteral<'minio'>
    }>
  ]
> = Type.Union([
  Type.Object({ type: Type.Literal('local') }),
  Type.Object({ type: Type.Literal('google-cloud') }),
  Type.Object({
    accesskey: Type.String(),
    bucket: Type.String(),
    endpoint: Type.String(),
    secretkey: Type.String(),
    type: Type.Literal('minio')
  })
])
export type StorageConfig = Static<typeof StorageConfigSchema>

export const RedisConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      auth: TOptional<TString>
      ca: TOptional<TString>
      enabled: TLiteral<true>
      host: TString
      port: TNumber
    }>
  ]
> = Type.Union([
  Type.Object({ enabled: Type.Literal(false) }),
  Type.Object({
    auth: Type.Optional(Type.String()),
    ca: Type.Optional(Type.String()),
    enabled: Type.Literal(true),
    host: Type.String(),
    port: Type.Number()
  })
])
export type RedisConfig = Static<typeof RedisConfigSchema>

export const SessionConfigSchema: TObject<{
  enabled: TBoolean
  secret: TString
}> = Type.Object({
  enabled: Type.Boolean({ default: true }),
  secret: Type.String({ default: 'session-scary-stuff' })
})
export type SessionConfig = Static<typeof SessionConfigSchema>

export const MonitoringConfigSchema: TObject<{
  enabled: TBoolean
  loki: TObject<{
    enabled: TBoolean
    host: TOptional<TString>
  }>
}> = Type.Object({
  enabled: Type.Boolean({ default: false }),
  loki: Type.Object({
    enabled: Type.Boolean({ default: false }),
    host: Type.Optional(Type.String())
  })
})
export type MonitoringConfig = Static<typeof MonitoringConfigSchema>

export const LoggingConfigSchema: TObject<{
  gcpfix: TBoolean
  level: TUnion<
    [
      TLiteral<'fatal'>,
      TLiteral<'error'>,
      TLiteral<'warn'>,
      TLiteral<'info'>,
      TLiteral<'debug'>,
      TLiteral<'trace'>,
      TLiteral<'silent'>
    ]
  >
  pretty: TBoolean
}> = Type.Object({
  gcpfix: Type.Boolean({ default: false }),
  level: Type.Union(
    [
      Type.Literal('fatal'),
      Type.Literal('error'),
      Type.Literal('warn'),
      Type.Literal('info'),
      Type.Literal('debug'),
      Type.Literal('trace'),
      Type.Literal('silent')
    ],
    { default: 'info' }
  ),
  pretty: Type.Boolean({ default: false })
})
export type LoggingConfig = Static<typeof LoggingConfigSchema>

export const ScraperConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      endpoint: TString
    }>
  ]
> = Type.Union([
  Type.Object({ enabled: Type.Literal(false) }),
  Type.Object({
    enabled: Type.Literal(true),
    endpoint: Type.String()
  })
])
export type ScraperConfig = Static<typeof ScraperConfigSchema>

export const UnstructuredConfigSchema: TUnion<
  [
    TObject<{
      enabled: TLiteral<false>
    }>,
    TObject<{
      enabled: TLiteral<true>
      endpoint: TString
    }>
  ]
> = Type.Union([
  Type.Object({ enabled: Type.Literal(false) }),
  Type.Object({
    enabled: Type.Literal(true),
    endpoint: Type.String()
  })
])
export type UnstructuredConfig = Static<typeof UnstructuredConfigSchema>

export const AuthConfigSchema: TObject<{
  firebase: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        clientEmail: TString
        enabled: TLiteral<true>
        privateKey: TString
        projectId: TString
      }>
    ]
  >
  local: TObject<{
    enabled: TBoolean
  }>
  twitter: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        callbackURL: TString
        consumerKey: TString
        consumerSecret: TString
        enabled: TLiteral<true>
      }>
    ]
  >
}> = Type.Object({
  firebase: Type.Union([
    Type.Object({ enabled: Type.Literal(false) }),
    Type.Object({
      clientEmail: Type.String(),
      enabled: Type.Literal(true),
      privateKey: Type.String(),
      projectId: Type.String()
    })
  ]),
  local: Type.Object({
    enabled: Type.Boolean({ default: false })
  }),
  twitter: Type.Union([
    Type.Object({ enabled: Type.Literal(false) }),
    Type.Object({
      callbackURL: Type.String(),
      consumerKey: Type.String(),
      consumerSecret: Type.String(),
      enabled: Type.Literal(true)
    })
  ])
})
export type AuthConfig = Static<typeof AuthConfigSchema>

export const ArchesConfigSchema: TObject<{
  auth: TObject<{
    firebase: TUnion<
      [
        TObject<{
          enabled: TLiteral<false>
        }>,
        TObject<{
          clientEmail: TString
          enabled: TLiteral<true>
          privateKey: TString
          projectId: TString
        }>
      ]
    >
    local: TObject<{
      enabled: TBoolean
    }>
    twitter: TUnion<
      [
        TObject<{
          enabled: TLiteral<false>
        }>,
        TObject<{
          callbackURL: TString
          consumerKey: TString
          consumerSecret: TString
          enabled: TLiteral<true>
        }>
      ]
    >
  }>
  billing: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        stripe: TObject<{
          token: TString
          whsec: TString
        }>
      }>
    ]
  >
  config: TObject<{
    validate: TBoolean
  }>
  database: TObject<{
    type: TUnion<
      [TLiteral<'postgres'>, TLiteral<'sqlite'>, TLiteral<'in-memory'>]
    >
    url: TString
  }>
  email: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        password: TString
        service: TString
        user: TString
      }>
    ]
  >
  embedding: TObject<{
    type: TUnion<[TLiteral<'openai'>, TLiteral<'ollama'>]>
  }>
  jwt: TObject<{
    expiration: TString
    secret: TString
  }>
  llm: TUnion<
    [
      TObject<{
        endpoint: TString
        token: TOptional<TString>
        type: TLiteral<'ollama'>
      }>,
      TObject<{
        endpoint: TOptional<TString>
        token: TString
        type: TLiteral<'openai'>
      }>
    ]
  >
  logging: TObject<{
    gcpfix: TBoolean
    level: TUnion<
      [
        TLiteral<'fatal'>,
        TLiteral<'error'>,
        TLiteral<'warn'>,
        TLiteral<'info'>,
        TLiteral<'debug'>,
        TLiteral<'trace'>,
        TLiteral<'silent'>
      ]
    >
    pretty: TBoolean
  }>
  monitoring: TObject<{
    enabled: TBoolean
    loki: TObject<{
      enabled: TBoolean
      host: TOptional<TString>
    }>
  }>
  platform: TObject<{
    enabled: TBoolean
    host: TString
  }>
  redis: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        auth: TOptional<TString>
        ca: TOptional<TString>
        enabled: TLiteral<true>
        host: TString
        port: TNumber
      }>
    ]
  >
  runpod: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        token: TString
      }>
    ]
  >
  scraper: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        endpoint: TString
      }>
    ]
  >
  server: TObject<{
    cors: TObject<{
      enabled: TBoolean
      origins: TString
    }>
    docs: TObject<{
      enabled: TBoolean
      export: TBoolean
    }>
    host: TString
    port: TNumber
  }>
  session: TObject<{
    enabled: TBoolean
    secret: TString
  }>
  speech: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        token: TString
      }>
    ]
  >
  storage: TUnion<
    [
      TObject<{
        type: TLiteral<'local'>
      }>,
      TObject<{
        type: TLiteral<'google-cloud'>
      }>,
      TObject<{
        accesskey: TString
        bucket: TString
        endpoint: TString
        secretkey: TString
        type: TLiteral<'minio'>
      }>
    ]
  >
  tls: TObject<{
    enabled: TBoolean
  }>
  unstructured: TUnion<
    [
      TObject<{
        enabled: TLiteral<false>
      }>,
      TObject<{
        enabled: TLiteral<true>
        endpoint: TString
      }>
    ]
  >
}> = Type.Object(
  {
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
  },
  {
    description: 'Arches AI configuration schema',
    title: 'Arches AI Configuration'
  }
)
export type ArchesConfig = Static<typeof ArchesConfigSchema>
