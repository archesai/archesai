'use client'

import { useCallback } from 'react'
import { useAtom } from 'jotai'

import {
  login as _login,
  logout as _logout,
  refresh as _refresh,
  register as _register,
  getOneUser
} from '@archesai/client'

import { authStatusAtom, defaultOrgnameAtom } from '#atoms/auth'

export const useAuth = () => {
  const [defaultOrgname, setDefaultOrgname] = useAtom(defaultOrgnameAtom)
  const [status, setStatus] = useAtom(authStatusAtom)

  const logout = useCallback(async () => {
    const response = await _logout({
      credentials: 'include',
      mode: 'cors'
    })
    console.log(response)
    setStatus('Unauthenticated')
  }, [setStatus])

  const getNewRefreshToken = useCallback(async (): Promise<void> => {
    if (status === 'Refreshing') {
      return
    }
    setStatus('Refreshing')
    const response = await _refresh({
      credentials: 'include',
      mode: 'cors'
    })
    if (response.status === 401) {
      throw new Error(response.data.errors[0]?.detail)
    }
    setStatus('Authenticated')
  }, [logout, setStatus, status])

  const authenticate = useCallback(async () => {
    try {
      let response = await getOneUser('me')
      if (response.status === 404) {
        await getNewRefreshToken()
        response = await getOneUser('me')
        if (response.status === 404) {
          throw new Error(response.data.errors[0]?.detail)
        }
      }
      setDefaultOrgname(response.data.data.attributes.orgname)
      setStatus('Authenticated')
    } catch (error: unknown) {
      console.error(error)
      return logout()
    }
  }, [logout, setStatus, setDefaultOrgname, getNewRefreshToken])

  const signInWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      const response = await _login(
        {
          email,
          password
        },
        {
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          mode: 'cors'
        }
      )

      if (response.status === 401) {
        throw new Error(response.data.errors[0]?.detail)
      }
    },
    []
  )

  const registerWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      const response = await _register(
        {
          email,
          password
        },
        {
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          mode: 'cors'
        }
      )

      if (response.status === 401) {
        throw new Error(response.data.errors[0]?.detail)
      }
    },
    []
  )

  return {
    authenticate,
    defaultOrgname,
    getNewRefreshToken,
    logout,
    registerWithEmailAndPassword,
    setStatus,
    signInWithEmailAndPassword,
    status
  }
}
