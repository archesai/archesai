import {
  ClassSerializerInterceptor,
  INestApplication,
  ValidationPipe,
} from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { Test, TestingModule } from "@nestjs/testing";
import { Logger, LoggerErrorInterceptor } from "nestjs-pino";
import request from "supertest";
import "tsconfig-paths/register";

import { ApiTokensService } from "../src/api-tokens/api-tokens.service";
import { RegisterDto } from "../src/auth/dto/register.dto";
import { TokenDto } from "../src/auth/dto/token.dto";
import { AppAuthGuard } from "../src/auth/guards/app-auth.guard";
import { EmailVerifiedGuard } from "../src/auth/guards/email-verified.guard";
import { OrganizationRoleGuard } from "../src/auth/guards/organization-role.guard";
import { RestrictedAPIKeyGuard } from "../src/auth/guards/restricted-api-key.guard";
import { AllExceptionsFilter } from "../src/common/filters/all-exceptions.filter";
import { ExcludeNullInterceptor } from "../src/common/interceptors/exclude-null.interceptor";
import { OrganizationEntity } from "../src/organizations/entities/organization.entity";
import { UserEntity } from "../src/users/entities/user.entity";
import { UsersService } from "../src/users/users.service";
import { AppModule } from "./../src/app.module"; // This enables path aliasing based on tsconfig.json
import { DeactivatedGuard } from "@/src/auth/guards/deactivated.guard";
import { RedisIoAdapter } from "@/src/common/adapters/redis-io.adapter";
import { ConfigService } from "@nestjs/config";
import RedisStore from "connect-redis";
import cookieParser from "cookie-parser";
import session from "express-session";
import { readFileSync } from "fs-extra";
import helmet from "helmet";
import passport from "passport";
import { createClient } from "redis";

export const createApp = async () => {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule],
  }).compile();
  const app = moduleFixture.createNestApplication();

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

  // Websocket Adapter
  const redisIoAdapter = new RedisIoAdapter(app, configService);
  await redisIoAdapter.connectToRedis();
  app.useWebSocketAdapter(redisIoAdapter);

  // Enable Shutdown Hooks
  app.enableShutdownHooks();

  return app;
};

export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

// Helper function to register a user and return the API token
export const registerUser = async (
  app: INestApplication,
  registerDto: RegisterDto
): Promise<TokenDto> => {
  const res = await request(app.getHttpServer())
    .post("/auth/register")
    .send(registerDto);
  expect(res.status).toBe(201);
  expect(res.type).toBe("application/json");
  expect(res).toSatisfyApiSpec();
  return res.body;
};

export const setEmailVerifiedByEmail = async (
  app: INestApplication,
  email: string
) => {
  const userService = app.get<UsersService>(UsersService);
  await userService.setEmailVerifiedByEmail(email);
};

// Helper function to get user data
export const getUser = async (
  app: INestApplication,
  accessToken: string
): Promise<UserEntity> => {
  const res = await request(app.getHttpServer())
    .get("/user")
    .set("Authorization", `Bearer ${accessToken}`);
  expect(res.status).toBe(200);
  expect(res.body.defaultOrgname).toBeTruthy();
  expect(res).toSatisfyApiSpec();
  return res.body;
};

// Helper function to check organization data
export const getOrganization = async (
  app: INestApplication,
  orgname: string,
  accessToken: string
): Promise<OrganizationEntity> => {
  const res = await request(app.getHttpServer())
    .get(`/organizations/${orgname}`)
    .set("Authorization", `Bearer ${accessToken}`);
  expect(res.status).toBe(200);
  expect(res).toSatisfyApiSpec();
  return res.body;
};

// Helper function to deactivate a user
export const deactivateUser = async (
  app: INestApplication,
  accessToken: string
) => {
  const res = await request(app.getHttpServer())
    .post("/user/deactivate")
    .set("Authorization", `Bearer ${accessToken}`);
  expect(res.status).toBe(201);
};
