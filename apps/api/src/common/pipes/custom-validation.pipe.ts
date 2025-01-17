import { BadRequestException, Injectable, ValidationPipe } from '@nestjs/common'

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
      stopAtFirstError: true,
      exceptionFactory: (errors) => {
        const messages = errors
          .filter((err) => err.constraints)
          .map((err) => Object.values(err.constraints as any))
          .flat()
          .join('; ')
        return new BadRequestException(messages)
      }
    })
  }
}
