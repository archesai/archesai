import { HttpModule } from "@nestjs/axios";
import { BullModule } from "@nestjs/bullmq";
import { Module } from "@nestjs/common";
import { ClassSerializerInterceptor, ValidationPipe } from "@nestjs/common";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { APP_FILTER, APP_GUARD, APP_INTERCEPTOR, APP_PIPE } from "@nestjs/core";
import { JwtModule } from "@nestjs/jwt";
import { MulterModule } from "@nestjs/platform-express";
import { ScheduleModule } from "@nestjs/schedule";
import { readFileSync } from "fs";
import Joi from "joi";
import { LoggerErrorInterceptor, LoggerModule } from "nestjs-pino";

import { ApiTokensModule } from "./api-tokens/api-tokens.module";
import { AudioModule } from "./audio/audio.module";
import { AuthModule } from "./auth/auth.module";
import { AppAuthGuard } from "./auth/guards/app-auth.guard";
import { DeactivatedGuard } from "./auth/guards/deactivated.guard";
import { EmailVerifiedGuard } from "./auth/guards/email-verified.guard";
import { OrganizationRoleGuard } from "./auth/guards/organization-role.guard";
import { RestrictedAPIKeyGuard } from "./auth/guards/restricted-api-key.guard";
import { BillingModule } from "./billing/billing.module";
import { AllExceptionsFilter } from "./common/filters/all-exceptions.filter";
import { ExcludeNullInterceptor } from "./common/interceptors/exclude-null.interceptor";
import { ContentModule } from "./content/content.module";
import { EmailModule } from "./email/email.module";
import { EmbeddingsModule } from "./embeddings/embeddings.module";
import { LabelsModule } from "./labels/labels.module";
import { LLMModule } from "./llm/llm.module";
import { MembersModule } from "./members/members.module";
import { OrganizationsModule } from "./organizations/organizations.module";
import { PipelinesModule } from "./pipelines/pipelines.module";
import { PrismaModule } from "./prisma/prisma.module";
import { RunpodModule } from "./runpod/runpod.module";
import { RunsModule } from "./runs/runs.module";
import { SpeechModule } from "./speech/speech.module";
import { StorageModule } from "./storage/storage.module";
import { ToolsModule } from "./tools/tools.module";
import { UsersModule } from "./users/users.module";
import { WebsocketsGateway } from "./websockets/websockets.gateway";
import { WebsocketsModule } from "./websockets/websockets.module";

@Module({
  controllers: [],
  imports: [
    PipelinesModule,
    LoggerModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (configService: ConfigService) => {
        const loggerConfig = {
          pinoHttp: {
            customProps: (req, res) => ({
              body: req?.body,
              context: "HTTP",
              origin: req?.headers?.origin,
              params: req?.params,
              path: req?.path?.split("?")[0],
              query: req?.query,
              statusCode: res?.statusCode,
            }),
            formatters:
              configService.get<string>("NODE_ENV") === "production"
                ? {
                    level(label: string) {
                      return { level: label, severity: label.toUpperCase() };
                    },
                  }
                : undefined,
            level:
              configService.get("NODE_ENV") === "test" ? "silent" : "debug",
            redact: {
              paths: ["req", "res"],
            },
            transport: {
              targets: [
                {
                  options: {
                    host:
                      configService.get<string>("LOKI_HOST") ||
                      "http://arches-loki:3100", // Loki service URL in Kubernetes
                    json: true,
                    labels: {
                      app: "archesai",
                      environment: "production",
                    },
                  },
                  target: "pino-loki",
                },
                {
                  options: {
                    colorize: true,
                    singleLine: true,
                  },
                  target: "pino-pretty",
                },
              ],
            },
          },
        };
        return loggerConfig;
      },
    }),
    AuthModule,
    UsersModule,
    OrganizationsModule,
    MembersModule,
    ConfigModule.forRoot({
      envFilePath: ["../.env"],
      ignoreEnvFile: process.env.NODE_ENV == "production",
      isGlobal: true,
      validationSchema: Joi.object({
        // CORS CONFIG
        ALLOWED_ORIGINS: Joi.string().required(),

        // DATABASE CONFIG
        DATABASE_URL: Joi.string().required(),
        EMAIL_PASSWORD: Joi.when("FEATURE_EMAIL", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),
        EMAIL_SERVICE: Joi.when("FEATURE_EMAIL", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),
        EMAIL_USER: Joi.when("FEATURE_EMAIL", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),

        // EMBEDDING CONFIG
        EMBEDDING_TYPE: Joi.string().valid("openai", "ollama").required(),

        // STRIPE CONFIG
        FEATURE_BILLING: Joi.boolean().required(),
        // EMAIL CONFIG
        FEATURE_EMAIL: Joi.boolean().required(),

        // JWT API TOKEN CONFIG
        JWT_API_TOKEN_EXPIRATION_TIME: Joi.string().required(),
        JWT_API_TOKEN_SECRET: Joi.string().required(),

        // LLM CONFIG
        LLM_TYPE: Joi.string().valid("openai", "ollama").required(),
        // LOADER CONFIG
        LOADER_ENDPOINT: Joi.string().required(),

        // GLOBAL CONFIG
        NODE_ENV: Joi.string().required(),

        OLLAMA_ENDPOINT: Joi.string().when("LLM_TYPE", {
          is: "ollama",
          otherwise: Joi.string().when("EMBEDDING_TYPE", {
            is: "ollama",
            otherwise: Joi.optional(),
            then: Joi.required(),
          }),
          then: Joi.required(),
        }),

        OPEN_AI_KEY: Joi.string().when("LLM_TYPE", {
          is: "openai",
          otherwise: Joi.string().when("EMBEDDING_TYPE", {
            is: "openai",
            otherwise: Joi.optional(),
            then: Joi.required(),
          }),
          then: Joi.required(),
        }),

        PORT: Joi.number().required(),

        // REDIS CONFIG
        REDIS_AUTH: Joi.string().required(),
        REDIS_CA_CERT_PATH: Joi.string().optional(),
        REDIS_HOST: Joi.string().required(),
        REDIS_PORT: Joi.number().required(),

        SERVER_HOST: Joi.string().required(),

        SESSION_SECRET: Joi.string().required(),

        // STORAGE TYPE
        STORAGE_TYPE: Joi.string()
          .valid("google-cloud", "local", "minio")
          .required(),

        STRIPE_PRIVATE_API_KEY: Joi.when("FEATURE_BILLING", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),

        STRIPE_WEBHOOK_SECRET: Joi.when("FEATURE_BILLING", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),
      }),
    }),
    BullModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        connection: {
          host: configService.get("REDIS_HOST"),
          password: configService.get("REDIS_AUTH"),
          port: Number(configService.get("REDIS_PORT")),
          tls: configService.get("REDIS_CA_CERT_PATH")
            ? {
                ca: readFileSync(configService.get("REDIS_CA_CERT_PATH")),
                rejectUnauthorized: false,
              }
            : undefined,
        },
      }),
    }),
    JwtModule,
    BillingModule,
    HttpModule,
    PrismaModule,
    EmailModule,
    MulterModule,
    ApiTokensModule,
    EmbeddingsModule,
    LLMModule,
    LabelsModule,
    StorageModule.forRoot(),
    WebsocketsModule,
    AudioModule,
    ScheduleModule.forRoot(),
    ContentModule,
    RunpodModule,
    SpeechModule,
    ToolsModule,
    RunsModule,
  ],
  providers: [
    WebsocketsGateway,
    {
      provide: APP_FILTER,
      useClass: AllExceptionsFilter,
    },
    {
      provide: APP_PIPE,
      useFactory: () => {
        return new ValidationPipe({
          forbidNonWhitelisted: true,
          forbidUnknownValues: true,
          transform: true,
          transformOptions: {
            enableImplicitConversion: true,
            exposeDefaultValues: true,
          },
          whitelist: true,
        });
      },
    },
    {
      provide: APP_GUARD,
      useClass: AppAuthGuard,
    },
    {
      provide: APP_GUARD,
      useClass: DeactivatedGuard,
    },
    {
      provide: APP_GUARD,
      useClass: EmailVerifiedGuard,
    },
    {
      provide: APP_GUARD,
      useClass: RestrictedAPIKeyGuard,
    },
    {
      provide: APP_GUARD,
      useClass: OrganizationRoleGuard,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ExcludeNullInterceptor,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ClassSerializerInterceptor,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: LoggerErrorInterceptor,
    },
  ],
})
export class AppModule {}

// {
//   target: "pino/file",
//   options: { destination: "/app-logs/app.log" },
// },
