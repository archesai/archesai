import type { Controller } from '@archesai/core'
import type { ApiTokenEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  API_TOKEN_ENTITY_KEY,
  ApiTokenEntitySchema,
  CreateApiTokenDtoSchema,
  UpdateApiTokenDtoSchema
} from '@archesai/schemas'

import type { ApiTokensService } from '#api-tokens/api-tokens.service'

export class ApiTokensController
  extends BaseController<ApiTokenEntity>
  implements Controller
{
  constructor(apiTokensService: ApiTokensService) {
    super(
      API_TOKEN_ENTITY_KEY,
      ApiTokenEntitySchema,
      CreateApiTokenDtoSchema,
      UpdateApiTokenDtoSchema,
      apiTokensService
    )
  }
}
