import type { JSX } from 'react'
import { useState } from 'react'

import type { RunEntity } from '#types/entities'

import {
  BanIcon,
  CheckCircle2Icon,
  ClockArrowUpIcon,
  Loader2Icon
} from '#components/custom/icons'
import { Button } from '#components/shadcn/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'

export const StatusTypeEnumButton = ({
  onClick,
  run
}: {
  onClick?: () => void
  run: RunEntity
  size?: 'lg' | 'sm'
}): JSX.Element => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false)

  const renderIcon = () => {
    switch (run.status) {
      case 'completed':
        return <CheckCircle2Icon className='text-green-500' />
      case 'failed':
        return <BanIcon className='text-destructive' />
      case 'pending':
        return <ClockArrowUpIcon className='text-orange-400' />
      case 'running':
        return <Loader2Icon className='animate-spin text-primary' />
      default:
        return null
    }
  }

  return (
    <Popover
      onOpenChange={setIsPopoverOpen}
      open={isPopoverOpen}
    >
      <PopoverTrigger asChild>
        <Button
          onClick={onClick}
          size='sm'
          variant='secondary'
        >
          {renderIcon()}
        </Button>
      </PopoverTrigger>
      <PopoverContent className='p-2 text-sm'>
        <div>
          <strong className='font-semibold'>Status:</strong> {run.status}
        </div>
        <div>
          <strong className='font-semibold'>Started:</strong>{' '}
          {run.startedAt && new Date(run.startedAt).toLocaleString()}
        </div>
        <div>
          <strong className='font-semibold'>Completed:</strong>{' '}
          {run.completedAt ? new Date(run.completedAt).toLocaleString() : 'N/A'}
        </div>
        {run.completedAt && (
          <div>
            <strong className='font-semibold'>Duration:</strong>{' '}
            {run.startedAt &&
              new Date(run.completedAt).getTime() -
                new Date(run.startedAt).getTime()}
          </div>
        )}
        {/* 
        <div>
          <strong className='font-semibold'>Progress:</strong>{' '}
          {Math.round(run.progress * 100)}% // FIXME
        </div> */}
        {run.error && (
          <div>
            <strong className='font-semibold'>Error:</strong> {run.error}
          </div>
        )}
      </PopoverContent>
    </Popover>
  )
}
