import { createServerFn } from '@tanstack/react-start'
import { getWebRequest } from '@tanstack/react-start/server'

import { getSession } from '@archesai/client'

export const getSessionServer = createServerFn({ method: 'GET' }).handler(
  async () => {
    const { headers } = getWebRequest()
    try {
      const { session, user } = (await getSession({
        credentials: 'include',
        headers
      })) as {
        session: null | { id: string; userId: string }
        user: null | { email: string; id: string; name: string }
      }
      if (user && session) {
        const sessionData = {
          session,
          user
        }
        return sessionData
      }
    } catch {
      /* empty */
    }
    return null
  }
)
