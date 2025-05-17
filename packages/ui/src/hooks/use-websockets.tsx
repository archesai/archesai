import { useEffect } from 'react'
import { io } from 'socket.io-client'

import { useAuth } from '#hooks/use-auth'

export const useWebsockets = ({
  overrideToken,
  queryClient
}: {
  overrideToken?: string
  queryClient?: {
    invalidateQueries: (options: {
      exact?: boolean
      queryKey: string[]
    }) => Promise<void>
  }
}) => {
  const { defaultOrgname: accessToken } = useAuth()

  useEffect(() => {
    if (accessToken) {
      const websocket = io(
        `ws${process.env.NEXT_PUBLIC_ARCHES_TLS_ENABLED ? 's' : ''}://${process.env.NEXT_PUBLIC_ARCHES_SERVER_HOST!}`,
        {
          auth: {
            token: overrideToken ?? accessToken
          },
          extraHeaders: {
            Authorization: `Bearer ${overrideToken ?? accessToken}`
          },
          reconnection: true,
          reconnectionAttempts: Infinity,
          reconnectionDelay: 1000,
          reconnectionDelayMax: 5000,
          transports: ['websocket'],
          withCredentials: true
        }
      )

      websocket.on('connect', () => {
        console.log('Connected to websocket')
      })

      websocket.on('ping', () => {
        console.log('Received ping')
        websocket.emit('pong')
      })

      websocket.on('update', async (event: { queryKey: string[] }) => {
        if (!queryClient) return
        await queryClient.invalidateQueries({
          queryKey: event.queryKey
        })
      })

      // websocket.on(
      //   'chat',
      //   (event: {
      //     content: ContentEntity
      //     labelId: string
      //     orgname: string
      //   }) => {
      //     streamContent(
      //       event.orgname,
      //       event.labelId,
      //       event.content,
      //       queryClient
      //     )
      //   }
      // )

      return () => {
        websocket.close()
      }
    }
    return
  }, [queryClient, accessToken, overrideToken])
}
