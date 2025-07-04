import ReactPlayer from 'react-player'

import type { ArtifactEntity } from '@archesai/domain'

// const ReactPlayer = dynamic(() => import('react-player'), { ssr: false })

export function ArtifactViewer({ content }: { content: ArtifactEntity }) {
  const { mimeType, text, url } = content

  let hoverContent = null

  if (!url) {
    return (
      <div className='flex h-full items-center justify-center'>
        <p>URL was not available</p>
      </div>
    )
  }

  if (mimeType.startsWith('video/') || mimeType.startsWith('audio/')) {
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
        src={url}
        width='100%'
      />
    )
  } else if (mimeType.startsWith('image/')) {
    hoverContent = (
      <image
        // className='h-full w-full object-contain'
        height={516}
        href={url}
        width={516}
      />
    )
  } else if (mimeType === 'application/pdf') {
    hoverContent = (
      <iframe
        className='h-full w-full'
        src={url}
        title='PDF Document'
      ></iframe>
    )
  } else if (mimeType.startsWith('text/')) {
    hoverContent = (
      <div className='flex h-full items-center justify-center p-4 text-center'>
        <p>{text}</p>
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
