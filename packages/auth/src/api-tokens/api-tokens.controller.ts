import type { Controller } from '@archesai/core'
import type { ApiTokenEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { API_TOKEN_ENTITY_KEY, ApiTokenEntitySchema } from '@archesai/domain'

import type { ApiTokensService } from '#api-tokens/api-tokens.service'

import { CreateApiTokenRequestSchema } from '#api-tokens/dto/create-api-token.req.dto'
import { UpdateApiTokenRequestSchema } from '#api-tokens/dto/update-api-token.req.dto'

export class ApiTokensController
  extends BaseController<ApiTokenEntity>
  implements Controller
{
  constructor(apiTokensService: ApiTokensService) {
    super(
      API_TOKEN_ENTITY_KEY,
      ApiTokenEntitySchema,
      CreateApiTokenRequestSchema,
      UpdateApiTokenRequestSchema,
      apiTokensService
    )
  }
}
