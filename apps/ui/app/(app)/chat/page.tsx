import Chat from '@/components/chat'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'
import { Suspense } from 'react'

export const metadata: Metadata = getMetadata('/chat')

export default function ChatPage() {
  return (
    <Suspense>
      <Chat />
    </Suspense>
  )
}
