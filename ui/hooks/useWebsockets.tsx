import { useQueryClient } from '@tanstack/react-query'
import { useEffect } from 'react'
import { io } from 'socket.io-client'

import { useAuth } from './use-auth'
import { useStreamChat } from './use-stream-chat'

export const useWebsockets = ({
  overrideToken
}: {
  overrideToken?: string
}) => {
  const { defaultOrgname: accessToken } = useAuth()
  const queryClient = useQueryClient()

  const { streamContent } = useStreamChat()

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

      websocket.on('connect', () => {
        console.debug('connected')
      })

      websocket.on('ping', () => {})

      websocket.on('update', async (event) => {
        await queryClient.invalidateQueries({
          queryKey: event.queryKey
        })
      })

      websocket.on('chat', (event) => {
        streamContent(event.orgname, event.labelId, event.content)
      })

      return () => {
        websocket.close()
      }
    }
  }, [queryClient, accessToken, overrideToken])
}
