/// <reference types="vite/client" />
import type { QueryClient } from '@tanstack/react-query'
import type { JSX } from 'react'

import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import type { GetOneSession200 } from '@archesai/client'

import { getGetOneSessionQueryKey } from '@archesai/client'
import { Toaster } from '@archesai/ui/components/shadcn/sonner'
import { LinkProvider } from '@archesai/ui/hooks/use-link'
import { seo } from '@archesai/ui/lib/seo'
import { ThemeProvider } from '@archesai/ui/providers/theme-provider'

import { DefaultCatchBoundary } from '#components/default-catch-boundary'
import NotFound from '#components/not-found'
import { SmartLink } from '#components/smart-links'
import getServerSession from '#lib/get-headers'
import globalsCss from '#styles/globals.css?url'

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient
  session: GetOneSession200 | null
}>()({
  beforeLoad: async ({ context }) => {
    const session = await context.queryClient.fetchQuery({
      queryFn: ({ signal }) => getServerSession({ signal }),
      queryKey: getGetOneSessionQueryKey()
    })
    return {
      session
    }
  },
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
}): JSX.Element {
  return (
    <html
      lang='en'
      suppressHydrationWarning
    >
      <head>
        <HeadContent />
      </head>
      <body className={`font-sans antialiased`}>
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
        {process.env.NODE_ENV === 'prod' && (
          <>
            <TanStackRouterDevtools position='bottom-right' />
            <ReactQueryDevtools buttonPosition='bottom-left' />
          </>
        )}
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
