/// <reference types="vite/client" />
import type { QueryClient } from '@tanstack/react-query'

import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import { Toaster } from '@archesai/ui/components/shadcn/sonner'
import { LinkProvider } from '@archesai/ui/hooks/use-link'
import { seo } from '@archesai/ui/lib/seo'
import { ThemeProvider } from '@archesai/ui/providers/theme-provider'

import { DefaultCatchBoundary } from '#components/default-catch-boundary'
import NotFound from '#components/not-found'
import { SmartLink } from '#components/smart-links'
import globalsCss from '../styles/globals.css?url'

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
}>()({
  component: RootComponent,
  errorComponent: (props) => {
    return (
      <RootDocument>
        <DefaultCatchBoundary {...props} />
      </RootDocument>
    )
  },
  head: () => ({
    links: [
      { href: globalsCss, rel: 'stylesheet' },
      {
        href: '/apple-touch-icon.png',
        rel: 'apple-touch-icon',
        sizes: '180x180'
      },
      {
        href: '/favicon-32x32.png',
        rel: 'icon',
        sizes: '32x32',
        type: 'image/png'
      },
      {
        href: '/favicon-16x16.png',
        rel: 'icon',
        sizes: '16x16',
        type: 'image/png'
      },
      { color: '#fffff', href: '/site.webmanifest', rel: 'manifest' },
      { href: '/favicon.ico', rel: 'icon' }
    ],
    meta: [
      { charSet: 'utf-8' },
      {
        content: 'width=device-width, initial-scale=1',
        name: 'viewport'
      },
      ...seo({
        description:
          'Arches AI is the perfect tool to explore documents using artificial intelligence. Simply upload your PDF and start asking questions to your personalized chatbot.',
        image: 'https://www.archesai.com/sc.png',
        title: 'Arches AI'
      })
    ]
  }),
  notFoundComponent: () => <NotFound />
})

export default function RootDocument({
  children
}: {
  children: React.ReactNode
}) {
  return (
    <html>
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
          <LinkProvider Link={SmartLink}>
            {children}
            <Toaster />
          </LinkProvider>
        </ThemeProvider>
        <TanStackRouterDevtools position='bottom-right' />
        <ReactQueryDevtools buttonPosition='bottom-left' />
        <Scripts />
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
