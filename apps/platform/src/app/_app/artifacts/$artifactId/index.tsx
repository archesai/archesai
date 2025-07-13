import { Suspense } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { ArtifactViewer } from '@archesai/ui/components/custom/artifact-viewer'
import { Card } from '@archesai/ui/components/shadcn/card'

import { ArtifactDetailsBody, ArtifactDetailsHeader } from './-details'

export const Route = createFileRoute('/_app/artifacts/$artifactId/')({
  component: ArtifactDetailsPage
})

export default function ArtifactDetailsPage() {
  const params = Route.useParams()

  return (
    <div className='flex h-full w-full gap-3'>
      {/*LEFT SIDE*/}
      <div className='flex w-1/2 flex-initial flex-col gap-3'>
        <Card>
          <Suspense>
            <ArtifactDetailsHeader artifactId={params.artifactId} />
          </Suspense>
          <Suspense>
            <ArtifactDetailsBody artifactId={params.artifactId} />
          </Suspense>
        </Card>
      </div>
      {/*RIGHT SIDE*/}
      <Card className='w-1/2 overflow-hidden'>
        <Suspense>
          <ArtifactViewer artifactId={params.artifactId} />
        </Suspense>
      </Card>
    </div>
  )
}
