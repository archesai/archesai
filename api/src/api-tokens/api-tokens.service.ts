import { Injectable } from '@nestjs/common'
import { JwtService } from '@nestjs/jwt'
import { v4 } from 'uuid'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { ApiTokenRepository } from './api-token.repository'
import { CreateApiTokenDto } from './dto/create-api-token.dto'
import { ApiTokenEntity, ApiTokenModel } from './entities/api-token.entity'
import { ArchesConfigService } from '../config/config.service'

@Injectable()
export class ApiTokensService extends BaseService<
  ApiTokenEntity,
  ApiTokenModel,
  ApiTokenRepository
> {
  constructor(
    private apiTokenRepository: ApiTokenRepository,
    private configService: ArchesConfigService,
    private jwtService: JwtService,
    private websocketsService: WebsocketsService
  ) {
    super(apiTokenRepository)
  }

  async create(
    data: CreateApiTokenDto & {
      username: string
      orgname: string
    }
  ) {
    const id = v4()
    const token = this.jwtService.sign(
      {
        domains: data.domains,
        id,
        orgname: data.orgname,
        role: data.role,
        username: data.username
      },
      {
        expiresIn: `${this.configService.get('jwt.expiration')}s`,
        secret: this.configService.get('jwt.secret')
      }
    )
    const key = '*********' + token.slice(-5)
    const apiToken = await this.apiTokenRepository.create({
      ...data,
      id,
      key
    })
    const entity = this.toEntity({
      ...apiToken,
      key: token
    })
    this.emitMutationEvent(entity)
    return entity
  }

  protected emitMutationEvent(entity: ApiTokenEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'api-tokens']
    })
  }

  protected toEntity(model: ApiTokenModel): ApiTokenEntity {
    return new ApiTokenEntity(model)
  }
}
