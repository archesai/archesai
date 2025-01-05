import { ArgumentMetadata, Injectable, Type } from '@nestjs/common'
import { CustomValidationPipe } from './custom-validation.pipe'

@Injectable()
export class AbstractValidationPipe extends CustomValidationPipe {
  constructor(
    private readonly targetTypes: {
      body?: Type
      query?: Type
      param?: Type
      custom?: Type
    } = {}
  ) {
    super()
  }

  async transform(value: any, metadata: ArgumentMetadata) {
    const targetType = this.targetTypes[metadata.type]
    if (!targetType) {
      return super.transform(value, metadata)
    }
    return super.transform(value, { ...metadata, metatype: targetType })
  }
}
