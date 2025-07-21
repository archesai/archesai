import { useCallback, useEffect, useState } from 'react'

type Listener = (view: ViewType) => void
type ViewType = 'grid' | 'table'

// Global state management
class ViewStateManager {
  get view() {
    return this._view
  }
  private _view: ViewType

  private listeners = new Set<Listener>()

  constructor() {
    this._view = this.getStoredView()

    // Listen for storage changes from other tabs
    if (typeof window !== 'undefined') {
      window.addEventListener('storage', this.handleStorageChange)
    }
  }

  setView(newView: ViewType) {
    this._view = newView
    if (typeof window !== 'undefined') {
      localStorage.setItem('view', newView)
    }
    this.notifyListeners()
  }

  subscribe(listener: Listener) {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  toggleView() {
    const newView = this._view === 'grid' ? 'table' : 'grid'
    this.setView(newView)
  }

  private getStoredView(): ViewType {
    if (typeof window === 'undefined') return 'grid'
    const storedView = localStorage.getItem('view')
    return storedView === 'table' ? 'table' : 'grid'
  }

  private handleStorageChange = (e: StorageEvent) => {
    if (e.key === 'view' && e.newValue) {
      const newView = e.newValue === 'table' ? 'table' : 'grid'
      this._view = newView
      this.notifyListeners()
    }
  }

  private notifyListeners() {
    this.listeners.forEach((listener) => {
      listener(this._view)
    })
  }
}

const viewStateManager = new ViewStateManager()

export const useToggleView = () => {
  const [view, setView] = useState<ViewType>(viewStateManager.view)

  useEffect(() => {
    const unsubscribe = viewStateManager.subscribe(setView)
    return unsubscribe
  }, [])

  useEffect(() => {
    const handleResize = () => {
      if (typeof window !== 'undefined' && window.innerWidth <= 768) {
        viewStateManager.setView('grid')
      }
    }

    if (typeof window !== 'undefined') {
      window.addEventListener('resize', handleResize)
      handleResize()
    }

    return () => {
      if (typeof window !== 'undefined') {
        window.removeEventListener('resize', handleResize)
      }
    }
  }, [])

  const setViewWrapper = useCallback((newView: ViewType) => {
    viewStateManager.setView(newView)
  }, [])

  const toggleView = useCallback(() => {
    viewStateManager.toggleView()
  }, [])

  return {
    setView: setViewWrapper,
    toggleView,
    view
  }
}
