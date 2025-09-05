import type { ThemeProviderProps } from 'next-themes'
import type { JSX } from 'react'

import { ThemeProvider as NextThemesProvider } from 'next-themes'

export const ThemeProvider = ({
  children,
  ...props
}: ThemeProviderProps): JSX.Element => {
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}
