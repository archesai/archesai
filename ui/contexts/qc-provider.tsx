'use client'

import { keepPreviousData, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'

export const QCProvider = ({ children }: any) => {
  const [client] = useState(
    new QueryClient({
      defaultOptions: {
        queries: {
          placeholderData: keepPreviousData
        }
      }
    })
  )

  return <QueryClientProvider client={client}>{children}</QueryClientProvider>
}
