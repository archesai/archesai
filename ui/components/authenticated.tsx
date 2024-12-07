'use client'

import { useAuth } from '@/hooks/use-auth'
import { useWebsockets } from '@/hooks/useWebsockets'
import { useEffect, useRef } from 'react'

export function Authenticated() {
  const { authenticate, getNewRefreshToken, status } = useAuth()
  const hasAuthenticated = useRef(false) // Track if authenticate() has been called

  useWebsockets({})

  useEffect(() => {
    if (status === 'Refreshing') {
      getNewRefreshToken()
    }

    if (status === 'Unauthenticated' && !hasAuthenticated.current) {
      ;(async () => {
        hasAuthenticated.current = true // Set the flag to prevent future calls
        await authenticate()
        hasAuthenticated.current = false // Reset the flag
      })()
    }
  }, [status])

  return null
}
