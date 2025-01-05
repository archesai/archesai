'use client'
import { baseUrl } from '@/generated/archesApiFetcher'
import { TokenDto } from '@/generated/archesApiSchemas'
import { useToast } from '@/hooks/use-toast'
import { GoogleAuthProvider, signInWithPopup } from 'firebase/auth'
import { useAtom } from 'jotai'
import { useCallback } from 'react'
import { auth } from '../lib/firebase'
import { authStatusAtom, defaultOrgnameAtom } from '../state/authState'
import { useRouter } from 'next/navigation'

export const useAuth = () => {
  const [defaultOrgname, setDefaultOrgname] = useAtom(defaultOrgnameAtom)
  const [status, setStatus] = useAtom(authStatusAtom)

  const router = useRouter()
  const { toast } = useToast()

  const logout = useCallback(async () => {
    console.log('baseUrl', baseUrl)
    const response = await fetch(baseUrl + '/auth/logout', {
      credentials: 'include',
      method: 'POST',
      mode: 'cors'
    })
    if (!response.ok) {
      const error = (await response.json()) as any
      console.error(error)
      toast({
        title: 'Failed to logout',
        description: error.message,
        variant: 'destructive'
      })
    }
    router.push('/')
    setStatus('Unauthenticated')
  }, [setStatus, router, toast])

  const getNewRefreshToken = useCallback(async () => {
    if (status === 'Refreshing') {
      return
    }
    setStatus('Refreshing')
    try {
      // Fetch new refresh token
      const response = await fetch(baseUrl + '/auth/refresh-token', {
        credentials: 'include',
        method: 'POST',
        mode: 'cors'
      })

      // If the response is not 201, throw an error
      if (!response.ok) {
        throw new Error(await response.text())
      }

      // Parse the response as a TokenDto
      const data = (await response.json()) as TokenDto
      setStatus('Authenticated')
      return data.accessToken
    } catch (error: any) {
      toast({
        title: 'Session expired. Please log in again.',
        description: error.message,
        variant: 'destructive'
      })
      await logout()
    }
  }, [logout, setStatus, status, toast])

  const fetchUser = useCallback(async () => {
    return fetch(baseUrl + '/user', {
      credentials: 'include', // Include cookies
      method: 'GET',
      mode: 'cors'
    })
  }, [])

  const authenticate = useCallback(async () => {
    try {
      // Attempt
      let response = await fetchUser()
      if (response.status === 401) {
        await getNewRefreshToken()
        response = await fetchUser()
        if (!response.ok) {
          return logout()
        }
      }
      const user = await response.json()
      setDefaultOrgname(user.defaultOrgname)
      setStatus('Authenticated')
    } catch (error: any) {
      console.error(error)
      toast({
        title: 'Authentication failed.',
        description: error.message,
        variant: 'destructive'
      })
      return logout()
    }
  }, [
    logout,
    setStatus,
    setDefaultOrgname,
    toast,
    getNewRefreshToken,
    fetchUser
  ])

  const signInWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
        const response = await fetch(baseUrl + '/auth/login', {
          body: JSON.stringify({ email, password }),
          credentials: 'include', // Include cookies
          headers: { 'Content-Type': 'application/json' },
          method: 'POST',
          mode: 'cors'
        })
        if (!response.ok) {
          const error = await response.json()
          throw new Error(error.message || 'Could not log in')
        }
        router.push('/playground')
      } catch (error: any) {
        toast({
          description: error.message,
          variant: 'destructive'
        })
        throw error
      }
    },
    [toast, router]
  )

  const registerWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
        const response = await fetch(baseUrl + '/auth/register', {
          body: JSON.stringify({ email, password }),
          credentials: 'include', // Include cookies
          headers: { 'Content-Type': 'application/json' },
          method: 'POST',
          mode: 'cors'
        })
        if (!response.ok) {
          const error = await response.json()
          throw new Error(error.message || 'Could not register user')
        }
        router.push('/playground')
      } catch (error: any) {
        console.error(error)
        toast({
          title: 'Registration failed.',
          description: error.message,
          variant: 'destructive'
        })
        throw error
      }
    },
    [toast, router]
  )

  const signInWithGoogle = useCallback(async () => {
    try {
      const provider = new GoogleAuthProvider()
      const credential = await signInWithPopup(auth, provider)
      const token = await credential.user.getIdToken()
      await fetch(baseUrl + '/auth/firebase/callback', {
        body: JSON.stringify({ accessToken: token }),
        credentials: 'include', // Include cookies
        headers: {
          'Content-Type': 'application/json'
        },
        method: 'POST',
        mode: 'cors'
      })

      router.push('/playground')
    } catch (error: any) {
      console.error(error)
      toast({
        title: 'Google sign-in failed.',
        description: error.message,
        variant: 'destructive'
      })
      await logout()
    }
  }, [logout, router, toast])

  return {
    authenticate,
    defaultOrgname,
    getNewRefreshToken,
    logout,
    registerWithEmailAndPassword,
    setStatus,
    signInWithEmailAndPassword,
    signInWithGoogle,
    status
  }
}
