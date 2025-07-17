import { GridIcon, ListIcon } from 'lucide-react'

import { Button } from '#components/shadcn/button'
import { useToggleView } from '#hooks/use-toggle-view'

export function ViewToggle() {
  const { setView, view } = useToggleView()
  return (
    <div className='hidden h-8 gap-2 md:flex'>
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === 'table' ? 'text-primary' : ''
        }`}
        onClick={() => {
          setView('table')
        }}
        variant={'outline'}
      >
        <ListIcon className='h-5 w-5' />
      </Button>
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === 'grid' ? 'text-primary' : ''
        }`}
        onClick={() => {
          setView('grid')
        }}
        variant={'outline'}
      >
        <GridIcon className='h-5 w-5' />
      </Button>
    </div>
  )
}
