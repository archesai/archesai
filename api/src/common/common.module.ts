import { ClassSerializerInterceptor, Module } from '@nestjs/common'
import { APP_FILTER, APP_INTERCEPTOR, APP_PIPE, Reflector } from '@nestjs/core'

import { AllExceptionsFilter } from './filters/all-exceptions.filter'
import { CustomValidationPipe } from './pipes/custom-validation.pipe'

@Module({
  providers: [
    {
      provide: APP_FILTER,
      useClass: AllExceptionsFilter
    },
    {
      provide: APP_PIPE,
      useFactory: () => new CustomValidationPipe()
    },
    {
      provide: APP_INTERCEPTOR,
      inject: [Reflector],
      useFactory: (reflector: Reflector) =>
        new ClassSerializerInterceptor(reflector, {
          excludeExtraneousValues: true,
          enableImplicitConversion: true
        })
    }
  ]
})
export class CommonModule {}
