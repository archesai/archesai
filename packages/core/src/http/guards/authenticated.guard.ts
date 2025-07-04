import type { ArchesApiRequest, ArchesApiResponse } from '@archesai/core'

import { UnauthorizedException } from '@archesai/core'

/**
 * Guard for authenticating with jwt, api-key-auth, or firebase-auth strategy.
 */
export function AuthenticatedGuard() {
  return async function (
    req: ArchesApiRequest,
    _reply: ArchesApiResponse
  ): Promise<void> {
    if (!req.user) {
      await Promise.reject(new UnauthorizedException())
    }
    await Promise.resolve()
  }
}
