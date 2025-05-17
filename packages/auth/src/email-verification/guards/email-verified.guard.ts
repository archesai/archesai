import type { ArchesApiRequest, CanActivate } from '@archesai/core'

import { ForbiddenException } from '@archesai/core'

/**
 * Guard for checking if the user has verified their e-mail.
 */
export class EmailVerifiedGuard implements CanActivate {
  public canActivate(request: ArchesApiRequest) {
    const { user } = request
    if (!user?.emailVerified) {
      throw new ForbiddenException(
        'You must verify your e-mail before using this feature.'
      )
    }

    return true
  }
}
