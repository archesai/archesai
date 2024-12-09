import { Suspense } from 'react'
import { Client } from './client'

export default function PlaygroundPage() {
  return (
    <Suspense fallback={<p>Loading feed...</p>}>
      <Client />
    </Suspense>
  )
}
