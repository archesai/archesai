import { useGetOneArtifactSuspense } from '@archesai/client'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { Badge } from '@archesai/ui/components/shadcn/badge'
import { Button } from '@archesai/ui/components/shadcn/button'
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

export const ArtifactDetailsHeader = ({
  artifactId
}: {
  artifactId: string
}) => {
  const {
    data: { data: artifact }
  } = useGetOneArtifactSuspense(artifactId)

  return (
    <CardHeader>
      <CardTitle className='flex items-center justify-between'>
        <div>{artifact.attributes.name}</div>
        <Button
          asChild
          size='sm'
          variant='outline'
        >
          <a
            href={artifact.attributes.url ?? ''}
            rel='noopener noreferrer'
            target='_blank'
          >
            Download Artifact
          </a>
        </Button>
      </CardTitle>
      <CardDescription>{artifact.attributes.description}</CardDescription>
    </CardHeader>
  )
}

export const ArtifactDetailsBody = ({ artifactId }: { artifactId: string }) => {
  const {
    data: { data: artifact }
  } = useGetOneArtifactSuspense(artifactId)

  return (
    <CardContent>
      <div className='flex items-center gap-2'>
        <Badge>{artifact.attributes.mimeType}</Badge>
        {artifact.attributes.createdAt && (
          <Badge>
            <Timestamp date={artifact.attributes.createdAt} />
          </Badge>
        )}
      </div>
    </CardContent>
  )
}
