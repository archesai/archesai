import { useNavigate } from '@tanstack/react-router'

import { useFindManyToolsSuspense } from '@archesai/client'
import { Card } from '@archesai/ui/components/shadcn/card'
import { cn, stringToColor } from '@archesai/ui/lib/utils'

import { toolBaseIcons } from '#lib/site-config'

export const ToolCards = () => {
  const navigate = useNavigate()
  const { data: toolsResponse } = useFindManyToolsSuspense()
  // if (isPending) {
  //   return (
  //     <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
  //       {new Array(6).map((_, index) => (
  //         <Card
  //           className='flex h-[150px] flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
  //           key={index}
  //         >
  //           <Skeleton className='mx-auto h-8 w-8' />
  //           <Skeleton className='mx-auto h-6 w-3/4' />
  //           <Skeleton className='mx-auto h-4 w-5/6' />
  //         </Card>
  //       ))}
  //     </div>
  //   )
  // }
  const tools = toolsResponse.data
  if (!tools.length) {
    return (
      <div className='flex h-[150px] items-center justify-center'>
        <p className='text-lg font-semibold'>No tools found</p>
      </div>
    )
  }
  return (
    <div className='grid grid-cols-1 gap-6 p-0 md:grid-cols-3'>
      {tools.map((tool, index) => {
        const Icon =
          toolBaseIcons[tool.attributes.toolBase as keyof typeof toolBaseIcons]
        const iconColor = stringToColor(tool.attributes.toolBase)
        return (
          <Card
            className='flex cursor-pointer flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
            onClick={async () => {
              await navigate({
                to: `/playground?selectedTool=${JSON.stringify(tool)}`
              })
            }}
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
            <div className='text-lg font-semibold'>{tool.attributes.name}</div>
            <div className='text-sm font-normal'>
              {' '}
              {tool.attributes.description}
            </div>
          </Card>
        )
      })}
    </div>
  )
}
