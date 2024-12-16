import {
  CanActivate,
  ExecutionContext,
  Injectable,
  Logger,
  UnauthorizedException
} from '@nestjs/common'
import { Reflector } from '@nestjs/core'
import { JwtService } from '@nestjs/jwt'
import { Request } from 'express'

@Injectable()
export class ApiTokenRestrictedDomainGuard implements CanActivate {
  private logger = new Logger(ApiTokenRestrictedDomainGuard.name)
  constructor(
    private reflector: Reflector,
    private jwtService: JwtService
  ) {}

  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest<Request>()

    // Extract the JWT token from the Authorization header
    const authHeader = request.headers.authorization
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return true
    }
    this.logger.debug(`Auth Header: ${authHeader}`)
    const token = authHeader.split(' ')[1]
    const payload = this.jwtService.decode(token)
    if (!payload || !payload.domains) {
      return true
    }
    const allowedDomains = payload.domains

    const clientIp = request.ip
    const origin = request.headers['origin'] || request.headers['referer']

    this.logger.debug(`Domains: ${allowedDomains}`)
    this.logger.debug(
      `Request Headers: ${JSON.stringify(request.headers, null, 2)}`
    )
    this.logger.debug(
      `Raw Headers: ${JSON.stringify(request.rawHeaders, null, 2)}`
    )
    this.logger.debug(`Client IP: ${clientIp}`)
    this.logger.debug(`Origin: ${origin}`)

    // Make sure if there was a domain in their header that this request is authorized for it
    if (allowedDomains && allowedDomains != '*') {
      if (
        !allowedDomains.split(',').includes(origin) &&
        !allowedDomains.split(',').includes(clientIp)
      ) {
        throw new UnauthorizedException('Unauthorized domain')
      }
    }

    return true
  }
}
