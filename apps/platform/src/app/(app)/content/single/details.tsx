'use client'

import { useSearchParams } from 'next/navigation'
import { format } from 'date-fns'

import { useGetOneContent } from '@archesai/client'
import { Badge } from '@archesai/ui/components/shadcn/badge'
import { Button } from '@archesai/ui/components/shadcn/button'
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

export const ContentDetailsHeader = () => {
  const searchParams = useSearchParams()
  const contentId = searchParams.get('contentId')!
  const { data } = useGetOneContent(contentId)

  if (!data) {
    return (
      <CardHeader>
        <CardTitle>Loading...</CardTitle>
      </CardHeader>
    )
  }

  if (data.status !== 200) {
    return (
      <CardHeader>
        <CardTitle>Error</CardTitle>
        <CardDescription>{data.data.errors[0]?.detail}</CardDescription>
      </CardHeader>
    )
  }

  const contentData = data.data.data
  return (
    <CardHeader>
      <CardTitle className='flex items-center justify-between'>
        <div>{contentData.attributes.name}</div>
        <Button
          asChild
          size='sm'
          variant='outline'
        >
          <a
            href={contentData.attributes.url ?? ''}
            rel='noopener noreferrer'
            target='_blank'
          >
            Download Content
          </a>
        </Button>
      </CardTitle>
      <CardDescription>{contentData.attributes.description}</CardDescription>
    </CardHeader>
  )
}

export const ContentDetailsBody = () => {
  const searchParams = useSearchParams()
  const contentId = searchParams.get('contentId')
  const { data: content } = useGetOneContent(contentId!)
  if (!content) {
    return (
      <CardContent>
        <div className='flex items-center gap-2'>Loading...</div>
      </CardContent>
    )
  }
  if (content.status !== 200) {
    return (
      <CardContent>
        <div className='flex items-center gap-2'>Error</div>
      </CardContent>
    )
  }
  const contentData = content.data.data
  return (
    <CardContent>
      <div className='flex items-center gap-2'>
        <Badge>{contentData.attributes.mimeType}</Badge>
        {contentData.attributes.createdAt && (
          <Badge>
            {format(new Date(contentData.attributes.createdAt), 'PPP')}
          </Badge>
        )}
      </div>
    </CardContent>
  )
}
