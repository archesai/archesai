import { createServerFn } from "@tanstack/react-start"
import { getWebRequest } from "@tanstack/react-start/server"

import type { GetOneSession200 } from "@archesai/client"

import { getOneSession } from "@archesai/client"

const getSessionServer = createServerFn({ method: "GET" }).handler(
  async (): Promise<GetOneSession200 | null> => {
    const { headers } = getWebRequest()
    try {
      const result = await getOneSession("current", {
        credentials: "include",
        headers
      })
      return result
    } catch {
      /* empty */
    }
    return null
  }
)

export default getSessionServer
