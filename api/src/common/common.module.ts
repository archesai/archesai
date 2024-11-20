import {
  ClassSerializerInterceptor,
  Module,
  ValidationPipe,
} from "@nestjs/common";
import { APP_FILTER, APP_INTERCEPTOR, APP_PIPE } from "@nestjs/core";

import { AllExceptionsFilter } from "./filters/all-exceptions.filter";
import { ExcludeNullInterceptor } from "./interceptors/exclude-null.interceptor";

@Module({
  providers: [
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
      provide: APP_INTERCEPTOR,
      useClass: ExcludeNullInterceptor,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ClassSerializerInterceptor,
    },
  ],
})
export class CommonModule {}
