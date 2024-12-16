import { Injectable, ValidationPipe } from '@nestjs/common'

@Injectable()
export class CustomValidationPipe extends ValidationPipe {
  constructor() {
    super({
      forbidNonWhitelisted: true,
      forbidUnknownValues: true,
      transform: true,
      transformOptions: {
        enableImplicitConversion: true,
        exposeDefaultValues: true
      },
      whitelist: true,
      skipMissingProperties: false,
      enableDebugMessages: true,
      stopAtFirstError: true
    })
  }
}
