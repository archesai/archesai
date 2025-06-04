import type { ArchesApiRequest, ArchesApiResponse } from '@archesai/core'

/**
 * Guard for authenticating with jwt, api-key-auth, or firebase-auth strategy.
 */
export function AuthenticatedGuard() {
  return async function (
    req: ArchesApiRequest,
    _reply: ArchesApiResponse
  ): Promise<void> {
    await new Promise<void>((resolve) => {
      // Check if the request has a user already
      if (req.user) {
        resolve()
        return
      }
    })
  }
}
