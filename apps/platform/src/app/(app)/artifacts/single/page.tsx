import { Suspense } from 'react'

// import { ArtifactViewer } from '@archesai/ui/components/custom/content-viewer'
import { Card } from '@archesai/ui/components/shadcn/card'

import { ArtifactDetailsBody, ArtifactDetailsHeader } from './details'

export default function ArtifactDetailsPage() {
  return (
    <div className='flex h-full w-full gap-3'>
      {/*LEFT SIDE*/}
      <div className='flex w-1/2 flex-initial flex-col gap-3'>
        <Card>
          <Suspense fallback={<p>Loading feed...</p>}>
            <ArtifactDetailsHeader />
          </Suspense>
          <Suspense fallback={<p>Loading feed...</p>}>
            <ArtifactDetailsBody />
          </Suspense>
        </Card>
      </div>
      {/*RIGHT SIDE*/}
      <Card className='w-1/2 overflow-hidden'>
        <Suspense fallback={<p>Loading feed...</p>}>
          {/* <ArtifactViewer /> FIXME */}
        </Suspense>
      </Card>
    </div>
  )
}
