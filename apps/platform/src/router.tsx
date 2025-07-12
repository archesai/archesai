import {
  MutationCache,
  notifyManager,
  QueryClient
} from '@tanstack/react-query'
import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routerWithQueryClient } from '@tanstack/react-router-with-query'

import { toast } from '@archesai/ui/components/shadcn/sonner'

import { DefaultCatchBoundary } from '#components/default-catch-boundary'
import NotFound from '#components/not-found'
import { routeTree } from './routeTree.gen'

export function createRouter() {
  if (typeof document !== 'undefined') {
    notifyManager.setScheduler(window.requestAnimationFrame)
  }

  const queryClient: QueryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnReconnect: () => !queryClient.isMutating(),
        refetchOnWindowFocus: false,
        retry: 0,
        staleTime: 1000 * 60 * 2 // 2 minutes
      }
    },
    mutationCache: new MutationCache({
      onError: (error) => {
        toast.error(error.message, { className: 'bg-red-500 text-white' })
      },
      onSettled: () => {
        if (queryClient.isMutating() === 1) {
          return queryClient.invalidateQueries()
        } else {
          return
        }
      }
    })
  })

  return routerWithQueryClient(
    createTanStackRouter({
      context: {
        queryClient
      },
      defaultErrorComponent: DefaultCatchBoundary,
      defaultNotFoundComponent: () => <NotFound />,
      defaultPendingComponent: () => {
        return (
          <div className='flex h-screen w-screen items-center justify-center bg-background'>
            <div className='h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent' />
          </div>
        )
      },
      defaultPreload: 'intent',
      // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#passing-all-loader-events-to-an-external-cache
      defaultPreloadStaleTime: 0,
      // defaultSsr: false,
      defaultStructuralSharing: true,
      routeTree,
      scrollRestoration: true
    }),
    queryClient
  )
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof createRouter>
  }
}
