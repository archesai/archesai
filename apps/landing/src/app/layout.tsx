import type { Metadata } from 'next'

import { Geist, Geist_Mono } from 'next/font/google'

import '../styles/globals.css'

const geistSans = Geist({
  subsets: ['latin'],
  variable: '--font-geist-sans'
})

const geistMono = Geist_Mono({
  subsets: ['latin'],
  variable: '--font-geist-mono'
})

export const fetchCache = 'default-cache'

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
} as Metadata

export default function RootLayout({
  children
}: {
  children: React.ReactNode
}) {
  return (
    <html
      lang='en'
      suppressHydrationWarning
    >
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  )
}
