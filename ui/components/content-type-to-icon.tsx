import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'
import { FileText, Image, Music, Video } from 'lucide-react'

export const ContentTypeToIcon = ({ contentType }: { contentType: string }) => {
  const getLabel = (contentType: string) => {
    switch (contentType) {
      case 'application/msword':
      case 'application/pdf':
      case 'application/vnd.ms-excel':
      case 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet':
      case 'application/vnd.openxmlformats-officedocument.wordprocessingml.document':
        return <FileText className='h-5 w-5 text-muted-foreground' />
      case 'audio/mp3':
      case 'audio/mpeg':
        return <Music className='h-5 w-5 text-muted-foreground' />
      case 'image/gif':
      case 'image/jpeg':
      case 'image/png':
      case 'image/svg+xml':
        return <Image className='h-5 w-5 text-muted-foreground' />

      case 'video/mp4':
      case 'video/mpeg':
        return <Video className='h-5 w-5 text-muted-foreground' />
      default:
        return <FileText className='h-5 w-5 text-muted-foreground' />
    }
  }

  return (
    <HoverCard>
      <HoverCardTrigger asChild>{getLabel(contentType)}</HoverCardTrigger>
      <HoverCardContent className='flex w-auto items-center gap-2 p-2'>
        {getLabel(contentType)}
        <p className='text-sm'>{contentType}</p>
      </HoverCardContent>
    </HoverCard>
  )
}
