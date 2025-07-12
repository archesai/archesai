import { randomUUID } from 'node:crypto'

import type { ConfigService, WebsocketsService } from '@archesai/core'
import type { ApiTokenEntity, BaseInsertion } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { API_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import type { ApiTokenRepository } from '#api-tokens/api-token.repository'
import type { JwtService } from '#jwt/jwt.service'

export class ApiTokensService extends BaseService<ApiTokenEntity> {
  private readonly apiTokenRepository: ApiTokenRepository
  private readonly configService: ConfigService
  private readonly jwtService: JwtService
  private readonly websocketsService: WebsocketsService

  constructor(
    apiTokenRepository: ApiTokenRepository,
    configService: ConfigService,
    jwtService: JwtService,
    websocketsService: WebsocketsService
  ) {
    super(apiTokenRepository)
    this.apiTokenRepository = apiTokenRepository
    this.configService = configService
    this.jwtService = jwtService
    this.websocketsService = websocketsService
  }

  public override async create(
    data: BaseInsertion<Omit<ApiTokenEntity, 'key'>>
  ): Promise<ApiTokenEntity> {
    const id = randomUUID()
    const token = this.jwtService.sign(
      {
        id,
        organizationId: data.organizationId,
        role: data.role
      },
      {
        expiresIn: Number.parseInt(this.configService.get('jwt.expiration'))
      }
    )
    const key = '*********' + token.slice(-5)
    const apiToken = await this.apiTokenRepository.create({
      ...data,
      id,
      key
    })
    this.emitMutationEvent(apiToken)
    return apiToken
  }

  protected emitMutationEvent(entity: ApiTokenEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, API_TOKEN_ENTITY_KEY]
    })
  }
}
