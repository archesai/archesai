// local-auth.guard.ts

import type {
  ArchesApiRequest,
  ArchesApiResponse,
  HttpInstance
} from '@archesai/core'

export function LocalAuthGuard(app: HttpInstance) {
  return async function (
    req: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<void> {
    await new Promise<void>((resolve, reject) => {
      const handler = req.passport.authenticate(
        ['local'],
        { session: true },
        async (authReq, _authRes, err, user) => {
          if (err) {
            reject(err)
            return
          }
          if (!user) {
            reject(new Error('Unauthorized'))
            return
          }

          try {
            await authReq.logIn(user)
            resolve()
          } catch (err) {
            reject(err as Error)
          }
        }
      )

      handler.call(app, req, reply)
    })
  }
}
