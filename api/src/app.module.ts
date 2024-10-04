import { HttpModule } from "@nestjs/axios";
import { BullModule } from "@nestjs/bull";
import { Module } from "@nestjs/common";
import { ConfigModule, ConfigService } from "@nestjs/config";
import { JwtModule } from "@nestjs/jwt";
import { MulterModule } from "@nestjs/platform-express";
import { ScheduleModule } from "@nestjs/schedule";
import { readFileSync } from "fs";
import * as Joi from "joi";
import { LoggerModule } from "nestjs-pino";

import { ApiTokensModule } from "./api-tokens/api-tokens.module";
import { ARTokensModule } from "./ar-tokens/ar-tokens.module";
import { AudioModule } from "./audio/audio.module";
import { AuthModule } from "./auth/auth.module";
import { ChatbotsModule } from "./chatbots/chatbots.module";
import { CompletionsModule } from "./completions/completions.module";
import { ContentModule } from "./content/content.module";
import { EmailModule } from "./email/email.module";
import { EmailChangeModule } from "./email-change/email-change.module";
import { EmailVerificationModule } from "./email-verification/email-verification.module";
import { EmbeddingsModule } from "./embeddings/embeddings.module";
import { FirebaseModule } from "./firebase/firebase.module";
import { JobsModule } from "./jobs/jobs.module";
import { MembersModule } from "./members/members.module";
import { MessagesModule } from "./messages/messages.module";
import { OrganizationsModule } from "./organizations/organizations.module";
import { PasswordResetModule } from "./password-reset/password-reset.module";
import { PrismaModule } from "./prisma/prisma.module";
import { RunpodModule } from "./runpod/runpod.module";
import { StorageModule } from "./storage/storage.module";
import { StripeModule } from "./stripe/stripe.module";
import { ThreadsModule } from "./threads/threads.module";
import { UsersModule } from "./users/users.module";
import { VectorDBModule } from "./vector-db/vector-db.module";
import { VectorRecordModule } from "./vector-records/vector-record.module";
import { WebsocketsGateway } from "./websockets/websockets.gateway";
import { WebsocketsModule } from "./websockets/websockets.module";

@Module({
  controllers: [],
  imports: [
    LoggerModule.forRoot({
      pinoHttp: {
        ...(process.env.NODE_ENV !== "production"
          ? {
              autoLogging: false, // This will disable automatic logging of HTTP requests
              transport: {
                target: "pino-pretty",
                // options: {
                //   singleLine: true,
                // },
              },
            }
          : {
              formatters: {
                level(label) {
                  return { level: label, severity: label.toUpperCase() };
                },
              },
            }),
        customProps: (req, res) => ({
          context: "HTTP",
          statusCode: res?.statusCode,
        }),
        redact: {
          paths: ["req.headers", "res.headers"],
        },
      },
    }),
    AuthModule,
    UsersModule,
    OrganizationsModule,
    MembersModule,
    ConfigModule.forRoot({
      envFilePath: ["../.env"],
      ignoreEnvFile: process.env.NODE_ENV == "production",
      validationSchema: Joi.object({
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
        // FIREBASE CONFIG
        FIREBASE_API_KEY: Joi.string().optional(),
        FRONTEND_HOST: Joi.string().required(),

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
        // PINECONE
        PINECONE_API_KEY: Joi.when("VECTOR_DB_TYPE", {
          is: "pinecone",
          otherwise: Joi.string().optional(),
          then: Joi.string().required(),
        }),
        PINECONE_INDEX: Joi.when("VECTOR_DB_TYPE", {
          is: "pinecone",
          otherwise: Joi.string().optional(),
          then: Joi.string().required(),
        }),
        PORT: Joi.number().required(),

        REDIS_AUTH: Joi.string().required(),

        REDIS_CA_CERT_PATH: Joi.string().optional(),

        // REDIS CONFIG
        REDIS_HOST: Joi.string().required(),
        REDIS_PORT: Joi.number().required(),
        SERVER_HOST: Joi.string().required(),

        // STORAGE TYPE
        STORAGE_TYPE: Joi.string()
          .valid("google-cloud", "local", "minio")
          .required(),
        STRIPE_API_CREDITS_PRICE_ID: Joi.when("FEATURE_BILLING", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),
        STRIPE_API_PRICE_ID: Joi.when("FEATURE_BILLING", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),
        STRIPE_BASIC_PRICE_ID: Joi.when("FEATURE_BILLING", {
          is: true,
          otherwise: Joi.string().forbidden(),
          then: Joi.string().required(),
        }),

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
        // STORAGE TYPE
        VECTOR_DB_TYPE: Joi.string().valid("pinecone", "pgvector").required(),
      }),
    }),
    BullModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => ({
        redis: {
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
    StripeModule,
    HttpModule,
    PrismaModule,
    EmailModule,
    MulterModule,
    ApiTokensModule,
    EmbeddingsModule,
    CompletionsModule,
    VectorDBModule,
    EmailVerificationModule,
    ThreadsModule,
    StorageModule.forRoot(),
    WebsocketsModule,
    ChatbotsModule,
    AudioModule,
    ScheduleModule.forRoot(),
    MessagesModule,
    JobsModule,
    ContentModule,
    RunpodModule,
    FirebaseModule,
    PasswordResetModule,
    EmailChangeModule,
    VectorRecordModule,
    ARTokensModule,
  ],
  providers: [WebsocketsGateway],
})
export class AppModule {}
