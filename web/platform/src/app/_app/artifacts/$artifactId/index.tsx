import type { JSX } from 'react'

import { Suspense } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { useGetOneArtifactSuspense } from '@archesai/client'
import { ArtifactViewer } from '@archesai/ui/components/custom/artifact-viewer'
import { Card } from '@archesai/ui/components/shadcn/card'

import {
  ArtifactDetailsBody,
  ArtifactDetailsHeader
} from '#app/_app/artifacts/$artifactId/-details'

export const Route = createFileRoute('/_app/artifacts/$artifactId/')({
  component: ArtifactDetailsPage
})

export default function ArtifactDetailsPage(): JSX.Element {
  const params = Route.useParams()
  const artifactId = params.artifactId

  return (
    <div className='flex h-full w-full gap-4'>
      {/*LEFT SIDE*/}
      <Card>
        <Suspense>
          <ArtifactDetailsHeader artifactId={artifactId} />
        </Suspense>
        <Suspense>
          <ArtifactDetailsBody artifactId={artifactId} />
        </Suspense>
      </Card>

      {/*RIGHT SIDE*/}
      <Card className='w-1/2 overflow-hidden'>
        <Suspense>
          <ArtifactViewerWrapper artifactId={artifactId} />
        </Suspense>
      </Card>
    </div>
  )
}

function ArtifactViewerWrapper({
  artifactId
}: {
  artifactId: string
}): JSX.Element {
  const {
    data: { data: artifact }
  } = useGetOneArtifactSuspense(artifactId)

  return <ArtifactViewer artifact={artifact} />
}
