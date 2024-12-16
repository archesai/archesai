'use client'
import { siteConfig } from '@/config/site'
import { useToolsControllerFindAll } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { useRouter } from 'next/navigation'

import { Card } from './ui/card'
import { Skeleton } from './ui/skeleton'
import { cn, stringToColor } from '@/lib/utils'

export const ToolCards = () => {
  const router = useRouter()
  const { defaultOrgname } = useAuth()
  const { data: tools, isPending } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  if (isPending) {
    return (
      <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
        {[...Array(6)].map((_, index) => (
          <Card
            className='flex h-[150px] flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
          >
            <Skeleton className='mx-auto h-8 w-8' />
            <Skeleton className='mx-auto h-6 w-3/4' />
            <Skeleton className='mx-auto h-4 w-5/6' />
          </Card>
        ))}
      </div>
    )
  }

  return (
    <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
      {tools?.results?.map((tool, index) => {
        const Icon = siteConfig.toolBaseIcons[tool.toolBase]
        const iconColor = stringToColor(tool.toolBase)
        return (
          <Card
            className='flex cursor-pointer flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
            onClick={() =>
              router.push(`/playground?selectedTool=${JSON.stringify(tool)}`)
            }
          >
            <Icon
              className={cn(
                'mx-auto h-8 w-8',
                iconColor.startsWith('text-') ? iconColor : ''
              )}
              style={{
                ...(iconColor.startsWith('#') ? { color: iconColor } : {})
              }}
            />
            <div className='text-lg font-semibold'>{tool.name}</div>
            <div className='text-sm font-normal'> {tool.description}</div>
          </Card>
        )
      })}
    </div>
  )
}
