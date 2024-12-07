'use client'
import { siteConfig } from '@/config/site'
import { useToolsControllerFindAll } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { useRouter } from 'next/navigation'

import { Button } from './ui/button'
import { Card } from './ui/card'

export const ToolCards = () => {
  const router = useRouter()
  const { defaultOrgname } = useAuth()
  const { data: tools } = useToolsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  return (
    <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
      {tools?.results?.map((tool, index) => {
        const Icon = siteConfig.toolBaseIcons[tool.toolBase]
        return (
          <Card
            className='flex flex-col justify-between gap-2 bg-sidebar p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
          >
            <Icon className='mx-auto h-8 w-8 text-primary/80' />
            <div className='text-lg font-semibold'>{tool.name}</div>
            <div className='text-sm font-normal'> {tool.description}</div>

            <Button
              className='mt-1 h-8'
              onClick={() => router.push(`/playground?selectedTool=${JSON.stringify(tool)}`)}
              variant={'outline'}
            >
              Select Tool
            </Button>
          </Card>
        )
      })}
    </div>
  )
}
