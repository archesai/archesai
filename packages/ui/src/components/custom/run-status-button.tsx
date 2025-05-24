'use client'

import { useState } from 'react'
import { Ban, CheckCircle2, ClockArrowUpIcon, Loader2Icon } from 'lucide-react'

import type { RunEntity } from '@archesai/domain'

import { Button } from '#components/shadcn/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'
import { cn } from '#lib/utils'

export const StatusTypeEnumButton = ({
  onClick,
  run
}: {
  onClick?: () => void
  run: RunEntity
  size?: 'lg' | 'sm'
}) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false)

  const renderIcon = () => {
    switch (run.status) {
      case 'COMPLETE':
        return <CheckCircle2 className='text-green-600' />
      case 'ERROR':
        return <Ban className='text-red-600' />
      case 'PROCESSING':
        return (
          <div className='flex items-center gap-2'>
            <Loader2Icon className='animate-spin text-primary' />
            <span>{(run.progress * 100).toFixed(0)}%</span>
          </div>
        )
      case 'QUEUED':
        return <ClockArrowUpIcon className='text-primary' />
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
          className={cn('flex items-center justify-between', '')}
          onClick={onClick}
          size='sm'
          variant='outline'
        >
          <div className='flex flex-1 items-center justify-start gap-1 truncate overflow-hidden'>
            {/* {Icon && <Icon className='text-blue-700' />} */}
          </div>
          <div className='ml-2 shrink-0'>{renderIcon()}</div>
        </Button>
      </PopoverTrigger>
      <PopoverContent className='overflow-auto p-4 text-sm'>
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

        <div>
          <strong className='font-semibold'>Progress:</strong>{' '}
          {Math.round(run.progress * 100)}%
        </div>
        {run.error && (
          <div>
            <strong className='font-semibold'>Error:</strong> {run.error}
          </div>
        )}
      </PopoverContent>
    </Popover>
  )
}
