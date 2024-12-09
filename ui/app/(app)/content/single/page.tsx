import { ContentViewer } from '@/components/content-viewer'
import { Card } from '@/components/ui/card'

import { Suspense } from 'react'
import { ContentDetailsBody, ContentDetailsHeader } from './details'

export default function ContentDetailsPage() {
  return (
    <div className='flex h-full w-full gap-3'>
      {/*LEFT SIDE*/}
      <div className='flex w-1/2 flex-initial flex-col gap-3'>
        <Card>
          <Suspense fallback={<p>Loading feed...</p>}>
            <ContentDetailsHeader />
          </Suspense>
          <Suspense fallback={<p>Loading feed...</p>}>
            <ContentDetailsBody />
          </Suspense>
        </Card>
      </div>
      {/*RIGHT SIDE*/}
      <Card className='w-1/2 overflow-hidden'>
        <Suspense fallback={<p>Loading feed...</p>}>
          <ContentViewer />
        </Suspense>
      </Card>
    </div>
  )
}
