import { useState } from 'react'
import {
  keepPreviousData,
  QueryClient,
  QueryClientProvider
} from '@tanstack/react-query'

export const QCProvider = ({ children }: { children: React.ReactNode }) => {
  const [client] = useState(
    new QueryClient({
      // defaultOptions: {
      //   dehydrate: {
      //     // include pending queries in dehydration
      //     shouldDehydrateQuery: (query) =>
      //       defaultShouldDehydrateQuery(query) ||
      //       query.state.status === 'pending'
      //   },
      //   queries: {
      //     staleTime: 60 * 1000
      //   }
      // }
      defaultOptions: {
        queries: {
          placeholderData: keepPreviousData
        }
      }
    })
  )

  return <QueryClientProvider client={client}>{children}</QueryClientProvider>
}
