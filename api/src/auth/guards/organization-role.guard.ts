import { UserEntity } from '@/src/users/entities/user.entity'
import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
  Logger,
  NotFoundException
} from '@nestjs/common'
import { Reflector } from '@nestjs/core'
import { Observable } from 'rxjs'

@Injectable()
export class MembershipGuard implements CanActivate {
  private readonly logger = new Logger(MembershipGuard.name)
  constructor(private reflector: Reflector) {}

  canActivate(
    context: ExecutionContext
  ): boolean | Observable<boolean> | Promise<boolean> {
    const { params, user } = context.switchToHttp().getRequest() as any
    const currentUser = user as UserEntity
    const orgname = params.orgname

    // Check for user and orgname, if they are not present this route is public and we skip. Alternatively, if orgname is not present, we skip.
    if (!orgname || !currentUser) {
      return true
    }

    // Check if user is a member of the organization
    const membership = currentUser.memberships.find(
      (val) => val.orgname == orgname
    )
    if (!membership) {
      this.logger.error(
        `User ${currentUser.username} is not a member of organization ${orgname}`
      )
      throw new NotFoundException()
    }

    // Check the roles that have access to this route
    const roles = this.reflector.getAllAndOverride<string[]>('roles', [
      context.getHandler(),
      context.getClass()
    ])
    if (!roles) {
      return true
    }

    // Check if user has the required role
    if (!roles.includes(membership.role)) {
      this.logger.error(
        `User ${currentUser.username} does not have the required role in organization ${orgname}`
      )
      throw new ForbiddenException(
        'You are not authorized to access this endpoint'
      )
    }

    return true
  }
}
