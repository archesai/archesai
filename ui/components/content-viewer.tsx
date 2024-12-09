'use client'
import dynamic from 'next/dynamic'
import Image from 'next/image'

import { useSearchParams } from 'next/navigation'
import { useAuth } from '@/hooks/use-auth'
import { fetchContentControllerFindOne } from '@/generated/archesApiComponents'
import { useSuspenseQuery } from '@tanstack/react-query'

const ReactPlayer = dynamic(() => import('react-player'), { ssr: false })

export function ContentViewer({ id }: { id?: string }) {
  const searchParams = useSearchParams()
  const contentId = searchParams?.get('contentId') || id
  const { defaultOrgname } = useAuth()

  const { data: content } = useSuspenseQuery({
    queryKey: ['organizations', defaultOrgname, 'content', contentId as string],
    queryFn: () =>
      fetchContentControllerFindOne({
        pathParams: {
          id: contentId as string,
          orgname: defaultOrgname
        }
      })
  })

  const { mimeType, text, url } = content || {}

  let hoverContent = null

  if (mimeType?.startsWith('video/') || mimeType?.startsWith('audio/')) {
    hoverContent = (
      <ReactPlayer
        config={{
          file: {
            attributes: {
              controlsList: 'nodownload'
            }
          }
        }}
        controls
        height='100%'
        url={url || ''}
        width='100%'
      />
    )
  } else if (mimeType?.startsWith('image/')) {
    hoverContent = (
      <Image
        alt={content?.description || ''}
        className='h-full w-full object-contain'
        height={516}
        src={url || ''}
        width={516}
      />
    )
  } else if (mimeType === 'application/pdf') {
    hoverContent = (
      <iframe
        className='h-full w-full'
        src={url || ''}
        title='PDF Document'
      ></iframe>
    )
  } else if (mimeType?.startsWith('text/')) {
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
