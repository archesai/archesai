import { applyDecorators, UseGuards } from '@nestjs/common'
import { ApiCookieAuth, ApiBearerAuth } from '@nestjs/swagger'
import { SetMetadata } from '@nestjs/common'
import { RoleTypeEnum } from '@/src/members/entities/member.entity'
import { AuthenticatedGuard } from '@/src/auth/guards/authenticated.guard'

export function Authenticated(roles: RoleTypeEnum[] = []) {
  SetMetadata('roles', roles)
  return applyDecorators(
    UseGuards(AuthenticatedGuard),
    ApiCookieAuth(),
    ApiBearerAuth()
  )
}
