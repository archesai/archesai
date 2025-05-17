import type {
  ArchesApiRequest,
  ArchesApiResponse,
  CanActivate,
  HttpInstance
} from '@archesai/core'

import { UnauthorizedException } from '@archesai/core'

/**
 * Guard for authenticating with jwt, api-key-auth, or firebase-auth strategy.
 */
export class AuthenticatedGuard implements CanActivate {
  private app: HttpInstance

  constructor(app: HttpInstance) {
    this.app = app
  }

  public async canActivate(
    request: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<boolean> {
    return new Promise((resolve, reject) => {
      const handler = request.passport.authenticate(
        ['jwt', 'api-key-auth'],
        { session: false },
        async (_req, _rep, err, user) => {
          if (err || !user) {
            reject(err ?? new UnauthorizedException())
          }
          resolve(true)
        }
      )
      handler.call(this.app, request, reply)
    })
  }
}
