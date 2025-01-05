'use client'

import { useAuth } from '@/hooks/use-auth'
import { useWebsockets } from '@/hooks/use-websockets'
import { useEffect } from 'react'

export function Authenticated() {
  const { authenticate, getNewRefreshToken, status } = useAuth()

  useWebsockets({})

  useEffect(() => {
    if (status === 'Loading') {
      ;(async () => {
        await authenticate()
      })()
    }
  }, [status, authenticate, getNewRefreshToken])

  return null
}
