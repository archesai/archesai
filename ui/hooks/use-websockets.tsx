import { useQueryClient } from '@tanstack/react-query'
import { useEffect } from 'react'
import { io } from 'socket.io-client'

import { useAuth } from './use-auth'
import { streamContent } from '@/lib/utils'

export const useWebsockets = ({
  overrideToken
}: {
  overrideToken?: string
}) => {
  const { defaultOrgname: accessToken } = useAuth()
  const queryClient = useQueryClient()

  useEffect(() => {
    if (accessToken) {
      const websocket = io(process.env.NEXT_PUBLIC_WEBSOCKET_URL as string, {
        auth: {
          token: overrideToken || accessToken
        },
        extraHeaders: {
          Authorization: `Bearer ${overrideToken || accessToken}`
        },
        reconnection: true,
        reconnectionAttempts: Infinity,
        reconnectionDelay: 1000,
        reconnectionDelayMax: 5000,
        transports: ['websocket'],
        withCredentials: true
      })

      websocket.on('connect', () => {})

      websocket.on('ping', () => {})

      websocket.on('update', async (event) => {
        await queryClient.invalidateQueries({
          queryKey: event.queryKey
        })
      })

      websocket.on('chat', (event) => {
        streamContent(event.orgname, event.labelId, event.content, queryClient)
      })

      return () => {
        websocket.close()
      }
    }
  }, [queryClient, accessToken, overrideToken])
}
