import { createServerFn } from '@tanstack/react-start'
import { getWebRequest } from '@tanstack/react-start/server'

import { getSession } from '@archesai/client'

export const getSessionServer = createServerFn({ method: 'GET' }).handler(
  async () => {
    const { headers } = getWebRequest()
    try {
      const { session, user } = await getSession({
        headers,
        credentials: 'include'
      })
      if (user && session) {
        const sessionData = {
          user,
          session
        }
        return sessionData
      }
    } catch {}
    return null
  }
)
