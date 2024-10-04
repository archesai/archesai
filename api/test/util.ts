import { ValidationPipe } from "@nestjs/common";
import { ClassSerializerInterceptor } from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { Test, TestingModule } from "@nestjs/testing";

import { ApiTokensService } from "../src/api-tokens/api-tokens.service";
import { AppAuthGuard } from "../src/auth/guards/app-auth.guard";
import { RestrictedAPIKeyGuard } from "../src/auth/guards/restricted-api-key.guard";
import { RolesGuard } from "../src/auth/guards/roles.guard";
import { AllExceptionsFilter } from "../src/common/all-exceptions.filter";
import { ExcludeNullInterceptor } from "../src/common/exclude-null.interceptor";
import { AppModule } from "./../src/app.module";

export const createApp = async () => {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule],
  }).compile();
  const app = moduleFixture.createNestApplication();

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
  app.useGlobalFilters(new AllExceptionsFilter());
  app.useGlobalInterceptors(
    new ClassSerializerInterceptor(app.get(Reflector)),
    new ExcludeNullInterceptor()
  );
  app.useGlobalGuards(
    new AppAuthGuard(app.get(Reflector)),
    new RestrictedAPIKeyGuard(app.get(Reflector), app.get(ApiTokensService)),
    new RolesGuard(app.get(Reflector))
  );

  app.enableShutdownHooks();

  return app;
};

export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
