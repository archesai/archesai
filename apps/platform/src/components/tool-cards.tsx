import { useNavigate } from '@tanstack/react-router'

import { useFindManyToolsSuspense } from '@archesai/client'
import { Card } from '@archesai/ui/components/shadcn/card'
import { cn, stringToColor } from '@archesai/ui/lib/utils'

import { toolBaseIcons } from '#lib/site-config'

export const ToolCards = () => {
  const navigate = useNavigate()
  const {
    data: { data: tools }
  } = useFindManyToolsSuspense()

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
        const Icon = toolBaseIcons['create-embeddings']
        const iconColor = stringToColor(tool.toolBase)
        return (
          <Card
            className='flex cursor-pointer flex-col justify-between gap-2 p-4 text-center transition-shadow hover:shadow-lg'
            key={index}
            onClick={async () => {
              await navigate({
                search: {
                  selectedTool: tool.id
                },
                to: `/playground`
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
            <div className='text-lg font-semibold'>{tool.name}</div>
            <div className='text-sm font-normal'>{tool.description}</div>
          </Card>
        )
      })}
    </div>
  )
}
