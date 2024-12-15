import Playground from '@/components/playground'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'
import { Suspense } from 'react'

export const metadata: Metadata = getMetadata('/playground')

export default function PlaygroundPage() {
  return (
    <Suspense>
      <Playground />
    </Suspense>
  )
}
