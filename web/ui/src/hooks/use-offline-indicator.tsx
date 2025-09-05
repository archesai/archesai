import { useEffect } from 'react'
import { onlineManager } from '@tanstack/react-query'
import { toast } from 'sonner'

export function useOfflineIndicator(): void {
  useEffect(() => {
    return onlineManager.subscribe(() => {
      if (onlineManager.isOnline()) {
        toast.success('online', {
          duration: 2000,
          id: 'ReactQuery'
        })
      } else {
        toast.error('offline', {
          duration: Infinity,
          id: 'ReactQuery'
        })
      }
    })
  }, [])
}
