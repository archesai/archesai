'use client'

import { useEffect } from 'react'

import { useAuth } from '#hooks/use-auth'
import { useWebsockets } from '#hooks/use-websockets'

export function Authenticated() {
  const { authenticate, getNewRefreshToken, status } = useAuth()

  useWebsockets({})

  useEffect(() => {
    if (status === 'Loading') {
      authenticate()
        .then()
        .catch((e: unknown) => {
          console.error(e)
        })
    }
  }, [status, authenticate, getNewRefreshToken])

  return null
}
