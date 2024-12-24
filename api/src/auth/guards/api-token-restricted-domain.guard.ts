import {
  CanActivate,
  ExecutionContext,
  Injectable,
  UnauthorizedException
} from '@nestjs/common'
import { Reflector } from '@nestjs/core'
import { JwtService } from '@nestjs/jwt'
import { Request } from 'express'

@Injectable()
export class ApiTokenRestrictedDomainGuard implements CanActivate {
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

    const token = authHeader.split(' ')[1]
    const payload = this.jwtService.decode(token)
    if (!payload || !payload.domains) {
      return true
    }
    const allowedDomains = payload.domains

    const clientIp = request.ip
    const origin = request.headers['origin'] || request.headers['referer']

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
