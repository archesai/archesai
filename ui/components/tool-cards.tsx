'use client'
import { siteConfig } from '@/config/site'
import { useToolsControllerFindAll } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { useRouter } from 'next/navigation'

import { Button } from './ui/button'
import { Card } from './ui/card'
import { Skeleton } from './ui/skeleton'

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
            className='flex h-[194px] flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
          >
            <Skeleton className='mx-auto h-8 w-8' />
            <Skeleton className='mx-auto h-6 w-3/4' />
            <Skeleton className='mx-auto h-4 w-5/6' />
            <Skeleton className='mt-1 h-8 w-full' />
          </Card>
        ))}
      </div>
    )
  }

  return (
    <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
      {tools?.results?.map((tool, index) => {
        const Icon = siteConfig.toolBaseIcons[tool.toolBase]
        return (
          <Card
            className='flex flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
          >
            <Icon className='mx-auto h-8 w-8 text-slate-700 dark:text-slate-500' />
            <div className='text-lg font-semibold'>{tool.name}</div>
            <div className='text-sm font-normal'> {tool.description}</div>

            <Button
              variant={'secondary'}
              className='mt-1 h-8'
              onClick={() =>
                router.push(`/playground?selectedTool=${JSON.stringify(tool)}`)
              }
            >
              Select Tool
            </Button>
          </Card>
        )
      })}
    </div>
  )
}
