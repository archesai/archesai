import type { ArchesApiRequest, CanActivate, Logger } from '@archesai/core'

import { ForbiddenException, UnauthorizedException } from '@archesai/core'

import type { MembersService } from '#members/members.service'

/**
 * Guard for checking user scopes.
 */
export class ScopesGuard implements CanActivate {
  private readonly logger: Logger
  private readonly membersService: MembersService

  constructor(membersService: MembersService, logger: Logger) {
    this.membersService = membersService
    this.logger = logger
  }

  public async canActivate(request: ArchesApiRequest): Promise<boolean> {
    const requiredScopes = ['scope1', 'scope2']

    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (!requiredScopes || requiredScopes.includes('public')) {
      this.logger.debug('Public scope; skipping permission check')
      return true
    }

    const { user } = request
    if (!user) {
      this.logger.debug('User is not authenticated')
      throw new UnauthorizedException('User is not authenticated')
    }

    this.logger.debug('Checking scopes for user', { requiredScopes, user })

    // Look up memberships for the user
    const { data: memberships } = await this.membersService.findMany({
      filter: {
        userId: {
          equals: user.id
        }
      }
    })

    // Extract user scopes from memberships
    const userScopes = memberships.map(
      (membership) => `organization:${membership.role}`
    )

    // Check if the user has any of the required scopes
    const hasPermission = requiredScopes.some((scope) =>
      userScopes.includes(scope)
    )

    if (!hasPermission) {
      this.logger.debug('User does not have required permissions', {
        requiredScopes,
        user,
        userScopes
      })
      throw new ForbiddenException('Insufficient permissions')
    }

    return true
  }
}
