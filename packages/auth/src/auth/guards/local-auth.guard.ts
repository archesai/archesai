import type {
  ArchesApiRequest,
  ArchesApiResponse,
  CanActivate,
  HttpInstance
} from '@archesai/core'

/**
 * Guard for authenticating with local strategy.
 */
export class LocalAuthGuard implements CanActivate {
  private app: HttpInstance

  constructor(app: HttpInstance) {
    this.app = app
  }

  public async canActivate(
    request: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<boolean> {
    return new Promise<boolean>((resolve) => {
      const handler = request.passport.authenticate(
        ['local'],
        { session: true },
        async (req, _rep, err, user) => {
          if (err || !user) {
            resolve(false)
          }
          await req.logIn(user)
          resolve(true)
        }
      )
      handler.call(this.app, request, reply)
    })
  }
}
