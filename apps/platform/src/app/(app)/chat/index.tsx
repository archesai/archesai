import { Suspense } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import Chat from '#components/chat'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/chat')

export const Route = createFileRoute('/(app)/chat/')({
  component: ChatPage
})

export default function ChatPage() {
  return (
    <Suspense>
      <Chat />
    </Suspense>
  )
}
