import { Toaster } from '@archesai/ui/components/shadcn/sonner'
import { NuqsAdapter } from '@archesai/ui/providers/nuqs-adapter'
import { QCProvider } from '@archesai/ui/providers/query-client-provider'
import { ThemeProvider } from '@archesai/ui/providers/theme-provider'

import '../styles/globals.css'

import {
  createRootRoute,
  HeadContent,
  Outlet,
  Scripts
} from '@tanstack/react-router'

// const fontSans = Geist({
//   subsets: ['latin'],
//   variable: '--font-sans'
// })

// const fontMono = Geist_Mono({
//   subsets: ['latin'],
//   variable: '--font-mono'
// })

export const Route = createRootRoute({
  component: RootLayout,
  head: () => ({
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

export default function RootLayout() {
  return (
    <html
      lang='en'
      suppressHydrationWarning
    >
      <head>
        <HeadContent />
      </head>
      <body
      // className={`${fontSans.variable} ${fontMono.variable} font-sans antialiased`}
      >
        <ThemeProvider
          attribute='class'
          defaultTheme='system'
          disableTransitionOnChange
          enableColorScheme
          enableSystem
        >
          <NuqsAdapter>
            <QCProvider>
              <Outlet />
              <Scripts />
            </QCProvider>
            <Toaster />
          </NuqsAdapter>
        </ThemeProvider>
      </body>
    </html>
  )
}
