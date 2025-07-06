/// <reference types="vite/client" />
import type { QueryClient } from '@tanstack/react-query'

// import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts
  //
} from '@tanstack/react-router'

import type { AuthContext } from '@archesai/ui/hooks/use-auth'

import { ActiveThemeProvider } from '@archesai/ui/components/custom/active-theme'
// import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import { Toaster } from '@archesai/ui/components/shadcn/sonner'
import { ThemeProvider } from '@archesai/ui/providers/theme-provider'

import globalsCss from '../styles/globals.css?url'

export const metadata = {
  description:
    'Arches AI is the perfect tool to explore documents using artificial intelligence. Simply upload your PDF and start asking questions to your personalized chatbot.',
  icons: {
    icon: '/icon.png'
  },
  openGraph: {
    description:
      'Arches AI is the perfect tool to explore documents using artificial intelligence. Simply upload your PDF and start asking questions to your personalized chatbot.',
    images: [
      {
        alt: 'Arches AI',
        height: 600,
        url: 'https://www.archesai.com/sc.png',
        width: 800
      }
    ],
    title: 'Arches AI',
    type: 'website',
    url: 'https://www.archesai.com/'
  },
  title: 'Arches AI',
  twitter: {
    card: 'summary_large_image',
    description:
      'Arches AI is the perfect tool to explore documents using artificial intelligence. Simply upload your PDF and start asking questions to your personalized chatbot.',
    images: ['https://www.archesai.com/sc.png'],
    title: 'Arches AI',
    url: 'https://www.archesai.com/'
  }
}

export const Route = createRootRouteWithContext<{
  authentication: AuthContext
  queryClient: QueryClient
}>()({
  component: RootComponent,
  head: () => ({
    links: [{ href: globalsCss, rel: 'stylesheet' }],
    meta: [
      { charSet: 'utf-8' },
      {
        content: 'width=device-width, initial-scale=1',
        name: 'viewport'
      },
      { title: 'TanStack Start Starter' }
    ]
  })
})

export default function RootDocument({
  children
}: {
  children: React.ReactNode
}) {
  return (
    <html
      lang='en'
      suppressHydrationWarning
    >
      <head>
        <HeadContent />
      </head>
      <body className={`overscroll-none font-sans antialiased`}>
        <ThemeProvider
          attribute='class'
          defaultTheme='system'
          disableTransitionOnChange
          enableColorScheme
          enableSystem
        >
          <ActiveThemeProvider>
            {children}
            {/* <TanStackRouterDevtools position='bottom-right' />
          <ReactQueryDevtools buttonPosition='bottom-left' /> */}
            <Scripts />
            <Toaster />
          </ActiveThemeProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}

function RootComponent() {
  return (
    <RootDocument>
      <Outlet />
    </RootDocument>
  )
}

// import { SiteFooter } from "@/components/site-footer"
// import { SiteHeader } from "@/components/site-header"

// export default function AppLayout({ children }: { children: React.ReactNode }) {
//   return (
//     <div className="bg-background relative z-10 flex min-h-svh flex-col">
//       <SiteHeader />
//       <main className="flex flex-1 flex-col">{children}</main>
//       <SiteFooter />
//     </div>
//   )
// }
