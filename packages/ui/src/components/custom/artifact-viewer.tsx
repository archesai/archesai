import { useSuspenseQuery } from '@tanstack/react-query'
import ReactPlayer from 'react-player'

import { getGetOneArtifactSuspenseQueryOptions } from '@archesai/client'

// const ReactPlayer = dynamic(() => import('react-player'), { ssr: false })

export function ArtifactViewer({ artifactId }: { artifactId: string }) {
  const {
    data: { data: artifact }
  } = useSuspenseQuery(getGetOneArtifactSuspenseQueryOptions(artifactId))
  let hoverContent: React.ReactNode = null
  if (
    artifact.attributes.mimeType.startsWith('video/') ||
    artifact.attributes.mimeType.startsWith('audio/')
  ) {
    hoverContent = (
      <ReactPlayer
        config={{
          html: {
            attributes: {
              controlsList: 'nodownload'
            }
          }
        }}
        controls
        height='100%'
        src={artifact.attributes.url ?? ''}
        width='100%'
      />
    )
  } else if (artifact.attributes.mimeType.startsWith('image/')) {
    hoverContent = (
      <image
        // className='h-full w-full object-contain'
        height={516}
        href={artifact.attributes.url}
        width={516}
      />
    )
  } else if (artifact.attributes.mimeType === 'application/pdf') {
    hoverContent = (
      <iframe
        className='h-full w-full'
        src={artifact.attributes.url}
        title='PDF Document'
      ></iframe>
    )
  } else if (artifact.attributes.mimeType.startsWith('text/')) {
    hoverContent = (
      <div className='flex h-full items-center justify-center p-4 text-center'>
        <p>{artifact.attributes.text}</p>
      </div>
    )
  } else {
    hoverContent = (
      <div className='flex h-full items-center justify-center'>
        <p>Cannot preview this content type. Please download to view.</p>
      </div>
    )
  }

  return hoverContent
}
