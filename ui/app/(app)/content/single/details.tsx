'use client'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { useContentControllerFindOne } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { format } from 'date-fns'
import { useSearchParams } from 'next/navigation'

export const ContentDetailsHeader = () => {
  const searchParams = useSearchParams()
  const contentId = searchParams.get('contentId')
  const { defaultOrgname } = useAuth()
  const { data: content } = useContentControllerFindOne(
    {
      pathParams: {
        id: contentId as string,
        orgname: defaultOrgname
      }
    },
    {
      enabled: !!defaultOrgname && !!contentId
    }
  )
  return (
    <CardHeader>
      <CardTitle className='flex items-center justify-between'>
        <div>{content?.name}</div>
        <Button
          asChild
          size='sm'
          variant='outline'
        >
          <a
            href={content?.url || ''}
            rel='noopener noreferrer'
            target='_blank'
          >
            Download Content
          </a>
        </Button>
      </CardTitle>
      <CardDescription>
        {content?.description || 'No Description'}
      </CardDescription>
    </CardHeader>
  )
}

export const ContentDetailsBody = () => {
  const searchParams = useSearchParams()
  const contentId = searchParams.get('contentId')
  const { defaultOrgname } = useAuth()
  const { data: content } = useContentControllerFindOne(
    {
      pathParams: {
        id: contentId as string,
        orgname: defaultOrgname
      }
    },
    {
      enabled: !!defaultOrgname && !!contentId
    }
  )
  return (
    <CardContent>
      <div className='flex items-center gap-2'>
        <Badge>{content?.mimeType}</Badge>
        {content?.createdAt && (
          <Badge>{format(new Date(content.createdAt), 'PPP')}</Badge>
        )}
      </div>
    </CardContent>
  )
}
