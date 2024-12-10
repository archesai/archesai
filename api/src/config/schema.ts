import Joi from 'joi'

export const validationSchema = Joi.object({
  // GLOBAL CONFIG
  NODE_ENV: Joi.string().required(),
  SERVER_HOST: Joi.string().required(),
  FRONTEND_HOST: Joi.string().required(),
  PORT: Joi.number().required(),
  ALLOWED_ORIGINS: Joi.string().required(),

  // DATABASE CONFIG
  DATABASE_URL: Joi.string().required(),

  // EMAIL CONFIG
  FEATURE_EMAIL: Joi.boolean().required(),
  EMAIL_USER: Joi.when('FEATURE_EMAIL', {
    is: true,
    otherwise: Joi.string().forbidden(),
    then: Joi.string().required()
  }),
  EMAIL_PASSWORD: Joi.when('FEATURE_EMAIL', {
    is: true,
    otherwise: Joi.string().forbidden(),
    then: Joi.string().required()
  }),
  EMAIL_SERVICE: Joi.when('FEATURE_EMAIL', {
    is: true,
    otherwise: Joi.string().forbidden(),
    then: Joi.string().required()
  }),

  // EMBEDDING CONFIG
  EMBEDDING_TYPE: Joi.string().valid('openai', 'ollama').required(),

  // JWT CONFIG
  JWT_API_TOKEN_EXPIRATION_TIME: Joi.string().required(),
  JWT_API_TOKEN_SECRET: Joi.string().required(),

  // BILLING CONFIG
  FEATURE_BILLING: Joi.boolean().required(),
  STRIPE_PRIVATE_API_KEY: Joi.when('FEATURE_BILLING', {
    is: true,
    otherwise: Joi.string().forbidden(),
    then: Joi.string().required()
  }),
  STRIPE_WEBHOOK_SECRET: Joi.when('FEATURE_BILLING', {
    is: true,
    otherwise: Joi.string().forbidden(),
    then: Joi.string().required()
  }),

  // LLM CONFIG
  LLM_TYPE: Joi.string().valid('openai', 'ollama').required(),
  LLM_ENDPOINT: Joi.string().when('LLM_TYPE', {
    is: 'ollama',
    otherwise: Joi.string().when('EMBEDDING_TYPE', {
      is: 'ollama',
      otherwise: Joi.optional(),
      then: Joi.required()
    }),
    then: Joi.required()
  }),
  LLM_API_KEY: Joi.string().when('LLM_TYPE', {
    is: 'openai',
    otherwise: Joi.string().when('EMBEDDING_TYPE', {
      is: 'openai',
      otherwise: Joi.optional(),
      then: Joi.required()
    }),
    then: Joi.required()
  }),

  // LOADER CONFIG
  LOADER_ENDPOINT: Joi.string().required(),

  // STORAGE TYPE
  STORAGE_TYPE: Joi.string().valid('google-cloud', 'local', 'minio').required(),

  // REDIS CONFIG
  REDIS_AUTH: Joi.string().required(),
  REDIS_CA_CERT_PATH: Joi.string().optional(),
  REDIS_HOST: Joi.string().required(),
  REDIS_PORT: Joi.number().required(),

  // SESSION CONFIG
  SESSION_SECRET: Joi.string().required(),

  // LOGGING CONFIG
  LOGGING_LEVEL: Joi.string()
    .valid('fatal', 'error', 'warn', 'info', 'debug', 'trace', 'silent')
    .required(),
  LOKI_HOST: Joi.string().optional()
})
