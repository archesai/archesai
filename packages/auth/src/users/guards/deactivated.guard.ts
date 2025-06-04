import type { ArchesApiRequest, CanActivate } from '@archesai/core'

import { ForbiddenException } from '@archesai/core'

/**
 * Guard that checks if the user is deactivated.
 */
export class DeactivatedGuard implements CanActivate {
  public canActivate(request: ArchesApiRequest) {
    const { user } = request
    if (user?.deactivated === true) {
      throw new ForbiddenException(
        'Your account has been deactivated. Please contact support.'
      )
    }

    return true
  }
}
