'use client'

import type { ThemeProviderProps } from 'next-themes'

import { useState } from 'react'
import {
  keepPreviousData,
  QueryClient,
  QueryClientProvider
} from '@tanstack/react-query'
import { ThemeProvider as NextThemesProvider } from 'next-themes'

export const ThemeProvider = ({ children, ...props }: ThemeProviderProps) => {
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}

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
