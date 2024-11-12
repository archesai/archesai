import { ClassSerializerInterceptor, ValidationPipe } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { NestFactory, Reflector } from "@nestjs/core";
import { DocumentBuilder, SwaggerModule } from "@nestjs/swagger";
import RedisStore from "connect-redis";
import cookieParser from "cookie-parser";
import session from "express-session";
import { readFileSync } from "fs-extra";
import helmet from "helmet";
import { Logger, LoggerErrorInterceptor } from "nestjs-pino";
import passport from "passport";
import { createClient } from "redis";

import { ApiTokensService } from "./api-tokens/api-tokens.service";
import { AppModule } from "./app.module";
import { AppAuthGuard } from "./auth/guards/app-auth.guard";
import { DeactivatedGuard } from "./auth/guards/deactivated.guard";
import { EmailVerifiedGuard } from "./auth/guards/email-verified.guard";
import { OrganizationRoleGuard } from "./auth/guards/organization-role.guard";
import { RestrictedAPIKeyGuard } from "./auth/guards/restricted-api-key.guard";
import { RedisIoAdapter } from "./common/adapters/redis-io.adapter";
import { AggregateFieldResult, Metadata } from "./common/dto/paginated.dto";
import {
  AggregateFieldQuery,
  FieldFieldQuery,
} from "./common/dto/search-query.dto";
import { AllExceptionsFilter } from "./common/filters/all-exceptions.filter";
import { ExcludeNullInterceptor } from "./common/interceptors/exclude-null.interceptor";

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    bufferLogs: true,
    rawBody: true,
  });
  const configService = app.get(ConfigService);

  //  Setup Logger
  app.useLogger(app.get(Logger));

  // Gloabl Filters and Interceptors
  app.useGlobalFilters(new AllExceptionsFilter());
  app.useGlobalInterceptors(
    new ExcludeNullInterceptor(),
    new ClassSerializerInterceptor(app.get(Reflector)),
    new LoggerErrorInterceptor()
  );

  // Global Pipes
  app.useGlobalPipes(
    new ValidationPipe({
      forbidNonWhitelisted: true,
      forbidUnknownValues: true,
      transform: true,
      transformOptions: {
        enableImplicitConversion: true,
        exposeDefaultValues: true,
      },
      whitelist: true,
    })
  );

  // Global Guards
  app.useGlobalGuards(
    new AppAuthGuard(app.get(Reflector)),
    new DeactivatedGuard(app.get(Reflector)),
    new EmailVerifiedGuard(app.get(Reflector)),
    new RestrictedAPIKeyGuard(app.get(Reflector), app.get(ApiTokensService)),
    new OrganizationRoleGuard(app.get(Reflector))
  );

  // CORS Configuration
  const allowedOrigins = configService
    .get<string>("ALLOWED_ORIGINS")
    .split(",");
  app.enableCors({
    allowedHeaders: ["Authorization", "Content-Type", "Accept"],
    credentials: true,
    origin: (origin, callback) => {
      if (!origin || allowedOrigins.includes(origin)) {
        callback(null, true);
      } else {
        callback(new Error("Not allowed by CORS"));
      }
    },
  });

  // Security Middlewares
  app.use(helmet());

  // Session Management
  const sessionSecret = configService.get<string>("SESSION_SECRET");
  if (!sessionSecret) {
    throw new Error("SESSION_SECRET is not defined");
  }
  const redisClient = createClient({
    password: configService.get("REDIS_AUTH"),
    url: `redis://${configService.get(
      "REDIS_HOST"
    )}:${configService.get("REDIS_PORT")}`,
    ...(configService.get("REDIS_CA_CERT_PATH")
      ? {
          socket: {
            ca: readFileSync(configService.get("REDIS_CA_CERT_PATH")),
            rejectUnauthorized: false,
            tls: true,
          },
        }
      : {}),
  });
  redisClient.on("error", (error) => {
    app.get(Logger).error("Redis client error: " + error);
  });
  redisClient.connect().catch(console.error);
  const redisStore = new RedisStore({
    client: redisClient,
  });

  app.use(cookieParser(sessionSecret));
  app.use(
    session({
      cookie: {
        httpOnly: true,
        maxAge: 24 * 60 * 60 * 1000,
        sameSite: "lax",
        secure: configService.get<string>("NODE_ENV") === "production",
      },
      resave: false,
      saveUninitialized: false,
      secret: sessionSecret,
      store: redisStore,
    })
  );

  // Initialize Passport
  app.use(passport.initialize());
  app.use(passport.session());

  // Swagger Setup
  if (configService.get<string>("NODE_ENV") !== "production") {
    const swaggerConfig = new DocumentBuilder()
      .setTitle("Arches AI API")
      .setDescription("The Arches AI API")
      .setVersion("v1")
      .addBearerAuth()
      .addServer(configService.get<string>("SERVER_HOST"))
      .build();

    const document = SwaggerModule.createDocument(app, swaggerConfig, {
      extraModels: [
        FieldFieldQuery,
        AggregateFieldQuery,
        AggregateFieldResult,
        Metadata,
      ],
    });

    SwaggerModule.setup("/", app, document, {
      customCss: ".swagger-ui .topbar { display: none }",
      swaggerOptions: {
        persistAuthorization: true,
        tagsSorter: "alpha",
      },
    });
  }

  // Websocket Adapter
  const redisIoAdapter = new RedisIoAdapter(app, configService);
  await redisIoAdapter.connectToRedis();
  app.useWebSocketAdapter(redisIoAdapter);

  // Enable Shutdown Hooks
  app.enableShutdownHooks();

  // Start listening for requests
  await app.listen(parseInt(configService.get<string>("PORT")) || 3000);
}
bootstrap();
