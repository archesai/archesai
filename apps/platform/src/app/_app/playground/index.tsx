import { createFileRoute } from '@tanstack/react-router'

import { Playground } from '#components/playground'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/playground')

export const Route = createFileRoute('/_app/playground/')({
  component: PlaygroundPage
})

export default function PlaygroundPage() {
  return <Playground />
}
