import { ClassSerializerInterceptor, ValidationPipe } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { NestFactory, Reflector } from "@nestjs/core";
import { DocumentBuilder, SwaggerModule } from "@nestjs/swagger";
import session from "express-session";
import { Logger, LoggerErrorInterceptor } from "nestjs-pino";

import { ApiTokensService } from "./api-tokens/api-tokens.service";
import {
  ApiTokenEntity,
  ChatbotsFieldItem,
} from "./api-tokens/entities/api-token.entity";
import { AppModule } from "./app.module";
import { AppAuthGuard } from "./auth/guards/app-auth.guard";
import { EmailVerifiedGuard } from "./auth/guards/email-verified.guard";
import { OrganizationRoleGuard } from "./auth/guards/organization-role.guard";
import { RestrictedAPIKeyGuard } from "./auth/guards/restricted-api-key.guard";
import { AllExceptionsFilter } from "./common/all-exceptions.filter";
import { ExcludeNullInterceptor } from "./common/exclude-null.interceptor";
import { _PaginatedDto } from "./common/paginated.dto";
import { RedisIoAdapter } from "./common/redis-io.adapter";
import { ThreadAggregates } from "./threads/dto/thread-aggregates.dto";

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    bufferLogs: true,
    rawBody: true,
  });
  app.useLogger(app.get(Logger));
  app.useGlobalFilters(new AllExceptionsFilter());
  app.useGlobalInterceptors(
    new ExcludeNullInterceptor(),
    new ClassSerializerInterceptor(app.get(Reflector)),
    new LoggerErrorInterceptor()
  );
  app.useGlobalPipes(
    new ValidationPipe({
      forbidUnknownValues: true,
      transform: true,
      transformOptions: {
        enableImplicitConversion: true,
        exposeDefaultValues: true,
      },
      whitelist: true,
    })
  );
  app.useGlobalGuards(
    new AppAuthGuard(app.get(Reflector)),
    new EmailVerifiedGuard(app.get(Reflector)),
    new RestrictedAPIKeyGuard(app.get(Reflector), app.get(ApiTokensService)),
    new OrganizationRoleGuard(app.get(Reflector))
  );
  app.enableCors({
    allowedHeaders: "Authorization, Content-Type, Accept",
    credentials: true,
    origin: "*",
  });
  const config = new DocumentBuilder()
    .setTitle("Arches AI API")
    .setDescription("The Arches AI API")
    .setVersion("v1")
    .addBearerAuth()
    .addServer(app.get(ConfigService).get("SERVER_HOST"))
    .build();

  const redisIoAdapter = new RedisIoAdapter(app, app.get(ConfigService));
  await redisIoAdapter.connectToRedis();

  app.useWebSocketAdapter(redisIoAdapter);

  app.use(
    session({
      cookie: {
        httpOnly: true, // Prevents client-side JS from accessing the cookie
        maxAge: 24 * 60 * 60 * 1000, // 1 day
        sameSite: "lax", // CSRF protection
        secure: false, // Set to true if using HTTPS
      },
      resave: false, // Do not save session if unmodified
      saveUninitialized: false, // Do not create session until something stored
      secret: process.env.SESSION_SECRET || "your_session_secret", // Use a strong secret in production
    })
  );

  const document = SwaggerModule.createDocument(app, config, {
    extraModels: [
      _PaginatedDto,
      ChatbotsFieldItem,
      ThreadAggregates,
      ApiTokenEntity,
    ],
  });

  SwaggerModule.setup("/", app, document, {
    customCss: ".swagger-ui .topbar { display: none }",
    swaggerOptions: {
      persistAuthorization: true,
      tagsSorter: "alpha",
    },
  });
  app.enableShutdownHooks();

  await app.listen(parseInt(process.env.PORT) || 3000);
}
bootstrap();
