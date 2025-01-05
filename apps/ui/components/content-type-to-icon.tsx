import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger
} from '@/components/ui/hover-card'
import { cn, stringToColor } from '@/lib/utils'
import {
  FileText,
  Image as ImageIcon,
  LetterText,
  Music,
  Video
} from 'lucide-react'

export const ContentTypeToIcon = ({ contentType }: { contentType: string }) => {
  const sharedClass = cn('h-5 w-5')
  const mediaType = contentType.split('/')[0] as string
  const color = stringToColor(contentType)
  const getLabel = (mediaType: string) => {
    switch (mediaType) {
      case 'application':
        return <FileText className={cn(sharedClass, color)} />
      case 'audio':
        return <Music className={cn(sharedClass, color)} />
      case 'image':
        return <ImageIcon className={cn(sharedClass, color)} />
      case 'video':
        return <Video className={cn(sharedClass, color)} />
      default:
        return (
          <LetterText
            className={cn(sharedClass)}
            style={{
              color: color
            }}
          />
        )
    }
  }

  return (
    <HoverCard>
      <HoverCardTrigger asChild>{getLabel(mediaType)}</HoverCardTrigger>
      <HoverCardContent className='flex w-auto items-center gap-2 p-2'>
        {getLabel(mediaType)}
        <p className='text-sm'>{contentType}</p>
      </HoverCardContent>
    </HoverCard>
  )
}
