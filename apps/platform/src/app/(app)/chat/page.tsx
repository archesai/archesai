import type { Metadata } from 'next'

import { Suspense } from 'react'

import Chat from '#components/chat'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/chat')

export default function ChatPage() {
  return (
    <Suspense>
      <Chat />
    </Suspense>
  )
}
