'use client'

import { useState } from 'react'
import {
  keepPreviousData,
  QueryClient,
  QueryClientProvider
} from '@tanstack/react-query'

export const QCProvider = ({ children }: { children: React.ReactNode }) => {
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
