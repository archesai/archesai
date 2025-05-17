import { useEffect, useState } from 'react'
import { useAtom } from 'jotai'

import { viewAtom } from '#atoms/view'

export const useToggleView = () => {
  const [view, setView] = useAtom(viewAtom)
  const [width, setWidth] = useState(0)

  const toggleView = () => {
    setView((prev) => (prev === 'grid' ? 'table' : 'grid'))
  }

  useEffect(() => {
    const handleResize = () => {
      setWidth(window.innerWidth)
    }

    window.addEventListener('resize', handleResize)

    if (width <= 768) {
      setView('grid')
    }

    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [width, setView])

  return {
    setView,
    toggleView,
    view
  }
}
