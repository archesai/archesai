'use client'

import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { useFilterItems } from '@/hooks/useFilterItems'
import { cn } from '@/lib/utils'
import { format } from 'date-fns'
import { Calendar as CalendarIcon } from 'lucide-react'
import * as React from 'react'

export function DatePickerWithRange() {
  const { range, setRange } = useFilterItems()

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          className={cn(
            'hidden h-8 justify-start gap-2 text-left font-normal md:flex',
            !range && 'text-muted-foreground'
          )}
          id='range'
          variant={'outline'}
        >
          <CalendarIcon className='h-5 w-5' />
          {range?.from ? (
            range.to ? (
              <>
                {format(range.from, 'LLL dd, y')} - {format(range.to, 'LLL dd, y')}
              </>
            ) : (
              format(range.from, 'LLL dd, y')
            )
          ) : (
            <span>Pick a date</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent align='start' className='w-auto p-0'>
        <Calendar
          defaultMonth={range?.from}
          initialFocus
          mode='range'
          numberOfMonths={2}
          onSelect={(range) => setRange(range as any)}
          selected={range}
        />
      </PopoverContent>
    </Popover>
  )
}
