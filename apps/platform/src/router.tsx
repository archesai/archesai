import { QueryClient } from '@tanstack/react-query'
import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routerWithQueryClient } from '@tanstack/react-router-with-query'

import { DefaultCatchBoundary } from '#components/default-catch-boundary'
import NotFound from '#components/not-found'
import { routeTree } from './routeTree.gen'

export function createRouter() {
  const queryClient = new QueryClient()

  return routerWithQueryClient(
    createTanStackRouter({
      context: {
        authentication: undefined!,
        queryClient
      },
      defaultErrorComponent: DefaultCatchBoundary,
      defaultNotFoundComponent: () => {
        return <NotFound />
      },
      defaultPreload: 'intent',
      // defaultPreloadStaleTime: 0,
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
