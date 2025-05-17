import type { Metadata } from 'next'

import { Suspense } from 'react'

import Playground from '#components/playground'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/playground')

export default function PlaygroundPage() {
  return (
    <Suspense>
      <Playground />
    </Suspense>
  )
}
