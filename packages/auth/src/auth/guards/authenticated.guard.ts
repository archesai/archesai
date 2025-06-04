import type {
  ArchesApiRequest,
  ArchesApiResponse,
  HttpInstance
} from '@archesai/core'

import { UnauthorizedException } from '@archesai/core'

/**
 * Guard for authenticating with jwt, api-key-auth, or firebase-auth strategy.
 */
export function AuthenticatedGuard(app: HttpInstance) {
  return async function (
    req: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<void> {
    await new Promise<void>((resolve, reject) => {
      const handler = req.passport.authenticate(
        ['local'],
        { session: true },
        async (_authReq, _authRes, err, user) => {
          if (err || !user) {
            reject(new UnauthorizedException())
            return
          }

          resolve()
        }
      )

      handler.call(app, req, reply)
    })
  }
}
