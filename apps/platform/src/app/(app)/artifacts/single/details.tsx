'use client'

import { useSearchParams } from 'next/navigation'

import { useGetOneContent } from '@archesai/client'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { Badge } from '@archesai/ui/components/shadcn/badge'
import { Button } from '@archesai/ui/components/shadcn/button'
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

export const ArtifactDetailsHeader = () => {
  const searchParams = useSearchParams()
  const artifactId = searchParams.get('artifactId')!
  const { data } = useGetOneContent(artifactId)

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

  const artifactData = data.data.data
  return (
    <CardHeader>
      <CardTitle className='flex items-center justify-between'>
        <div>{artifactData.attributes.name}</div>
        <Button
          asChild
          size='sm'
          variant='outline'
        >
          <a
            href={artifactData.attributes.url ?? ''}
            rel='noopener noreferrer'
            target='_blank'
          >
            Download Artifact
          </a>
        </Button>
      </CardTitle>
      <CardDescription>{artifactData.attributes.description}</CardDescription>
    </CardHeader>
  )
}

export const ArtifactDetailsBody = () => {
  const searchParams = useSearchParams()
  const artifactId = searchParams.get('artifactId')
  const { data: content } = useGetOneContent(artifactId!)
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
  const artifactData = content.data.data
  return (
    <CardContent>
      <div className='flex items-center gap-2'>
        <Badge>{artifactData.attributes.mimeType}</Badge>
        {artifactData.attributes.createdAt && (
          <Badge>
            <Timestamp date={artifactData.attributes.createdAt} />
          </Badge>
        )}
      </div>
    </CardContent>
  )
}
