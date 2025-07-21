import { useEffect, useState } from 'react'

export const useToggleView = () => {
  const [view, setView] = useState<'grid' | 'table'>('grid')
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
