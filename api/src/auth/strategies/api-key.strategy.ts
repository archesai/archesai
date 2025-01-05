import { UserEntity } from '@/src/users/entities/user.entity'
import { Injectable, UnauthorizedException } from '@nestjs/common'
import { Logger } from '@nestjs/common'
import { PassportStrategy } from '@nestjs/passport'
import { ExtractJwt, Strategy } from 'passport-jwt'

import { ApiTokensService } from '../../api-tokens/api-tokens.service'
import { UsersService } from '../../users/users.service'
import { OperatorEnum } from '@/src/common/dto/search-query.dto'
import { ArchesConfigService } from '@/src/config/config.service'

@Injectable()
export class ApiKeyStrategy extends PassportStrategy(Strategy, 'api-key-auth') {
  private readonly logger: Logger = new Logger(ApiKeyStrategy.name)

  constructor(
    private configService: ArchesConfigService,
    private usersService: UsersService,
    private apiTokensService: ApiTokensService
  ) {
    super({
      ignoreExpiration: false,
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: configService.get('jwt.secret')
    })
  }

  async validate(payload: any): Promise<UserEntity> {
    this.logger.debug(`Validating API Key: ${payload.id}`)
    const { id, orgname, role, username } = payload
    const user = await this.usersService.findOneByUsername(username)
    this.logger.debug(
      `Found user valid for API Key: ${user.username}, checking memberships`
    )

    user.memberships = user.memberships.filter((m) => m.orgname == orgname)
    if (!user.memberships.length) {
      throw new UnauthorizedException()
    }

    const tokens = await this.apiTokensService.findAll({
      filters: [
        {
          field: 'orgname',
          operator: OperatorEnum.EQUALS,
          value: orgname
        }
      ]
    })
    if (!tokens || !tokens.results.find((t) => t.id == id)) {
      throw new UnauthorizedException()
    }
    user.memberships[0].role = role

    return user
  }
}
