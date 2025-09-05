import { useEffect, useState } from 'react'

type ViewType = 'grid' | 'table'

// Cookie utility functions
const getCookie = (name: string): null | string => {
  if (typeof document === 'undefined') return null

  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() ?? null
  }
  return null
}

const setCookie = (name: string, value: string, days = 30) => {
  if (typeof document === 'undefined') return

  const expires = new Date()
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000)
  document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`
}

export const useToggleView = ({
  defaultView = 'table'
}: { defaultView?: ViewType } = {}): {
  setView: (newView: ViewType) => void
  toggleView: () => void
  view: ViewType
} => {
  // Initialize from cookie or use default
  const getInitialView = (): ViewType => {
    const savedView = getCookie('viewType') as null | ViewType
    return savedView ?? defaultView
  }

  const [view, setView] = useState<ViewType>(getInitialView)

  // Handle responsive behavior
  useEffect(() => {
    const handleResize = () => {
      if (typeof window !== 'undefined' && window.innerWidth <= 768) {
        const newView = 'grid'
        setView(newView)
        setCookie('viewType', newView)
      }
    }

    if (typeof window !== 'undefined') {
      window.addEventListener('resize', handleResize)
      handleResize() // Check on mount
    }

    return () => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('resize', handleResize)
      }
    }
  }, [])

  const setViewWrapper = (newView: ViewType) => {
    setView(newView)
    setCookie('viewType', newView)
  }

  const toggleView = () => {
    const newView = view === 'grid' ? 'table' : 'grid'
    setView(newView)
    setCookie('viewType', newView)
  }

  return {
    setView: setViewWrapper,
    toggleView,
    view
  }
}
