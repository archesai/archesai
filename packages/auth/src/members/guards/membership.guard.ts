import type { ArchesApiRequest, CanActivate } from '@archesai/core'

import { Logger } from '@archesai/core'

import type { MembersService } from '#members/members.service'

/**
 * Guard for checking user membership.
 */
export class MembershipGuard implements CanActivate {
  private readonly logger = new Logger(MembershipGuard.name)
  private readonly membersService: MembersService

  constructor(membersService: MembersService) {
    this.membersService = membersService
  }

  public async canActivate(
    request: ArchesApiRequest & { params: { orgname?: string } }
  ): Promise<boolean> {
    const orgname = request.params.orgname
    const user = request.user

    // If no organization is specified or no user is attached, allow the request.
    if (!orgname || !user) {
      return true
    }

    this.logger.debug('Checking membership for user', {
      orgname,
      path: request.url,
      user
    })

    // Allow public route for joining an organization
    if (/^\/organizations\/[^/]+\/members\/join$/.test(request.url)) {
      this.logger.debug('Public join route; skipping membership check')
      return true
    }

    // Look up memberships for the user
    const { data: memberships } = await this.membersService.findMany({
      filter: {
        userId: {
          equals: user.id
        }
      }
    })

    // Find membership that matches the orgname
    const membership = memberships.find((m) => m.orgname === orgname)
    if (!membership) {
      this.logger.debug('user is not a member', { orgname, user })
      return false
    }

    // Check allowed roles using the rolesGetter helper
    const allowedRoles = ['OWNER', 'ADMIN'] as string[]
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
    if (allowedRoles && allowedRoles.length > 0) {
      if (!allowedRoles.includes(membership.role)) {
        this.logger.debug('User does not have required role', {
          orgname,
          requiredRoles: allowedRoles,
          user,
          userRole: membership.role
        })
        return false
      }
    }

    return true
  }
}
