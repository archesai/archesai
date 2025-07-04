import { useGetOneArtifact } from '@archesai/client'
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
  const artifactId = ''
  const { data, error } = useGetOneArtifact(artifactId)

  if (!data) {
    return (
      <CardHeader>
        <CardTitle>Loading...</CardTitle>
      </CardHeader>
    )
  }

  if (error) {
    return (
      <CardHeader>
        <CardTitle>Error</CardTitle>
        <CardDescription>{error.errors[0]?.detail}</CardDescription>
      </CardHeader>
    )
  }

  const artifactData = data.data
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
  const artifactId = ''
  const { data: artifact, error } = useGetOneArtifact(artifactId)
  if (!artifact) {
    return (
      <CardContent>
        <div className='flex items-center gap-2'>Loading...</div>
      </CardContent>
    )
  }
  if (error) {
    return (
      <CardContent>
        <div className='flex items-center gap-2'>Error</div>
      </CardContent>
    )
  }
  const artifactData = artifact.data
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
