'use client'

import { buttonVariants } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ChevronLeftIcon, ChevronRightIcon } from '@radix-ui/react-icons'
import * as React from 'react'
import { DayPicker } from 'react-day-picker'

export type CalendarProps = React.ComponentProps<typeof DayPicker>

function Calendar({ className, classNames, showOutsideDays = true, ...props }: CalendarProps) {
  return (
    <DayPicker
      className={cn('p-3', className)}
      classNames={{
        button_next: cn(
          buttonVariants({
            className: 'absolute right-0 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100',
            variant: 'outline'
          })
        ),
        button_previous: cn(
          buttonVariants({
            className: 'absolute left-0 h-7 w-7 bg-transparent p-0 opacity-50 hover:opacity-100',
            variant: 'outline'
          })
        ),
        caption: 'flex justify-center pt-1 relative items-center',
        caption_label: 'text-sm font-medium truncate',
        day: 'p-0 size-8 text-sm flex-1 flex items-center justify-center has-[button]:hover:!bg-accent rounded-md has-[button]:hover:aria-selected:!bg-primary has-[button]:hover:text-accent-foreground has-[button]:hover:aria-selected:text-primary-foreground',
        day_button: cn(
          buttonVariants({ variant: 'ghost' }),
          'size-8 p-0 font-normal transition-none hover:bg-transparent hover:text-inherit aria-selected:opacity-100'
        ),
        disabled: 'text-muted-foreground opacity-50',
        hidden: 'invisible',
        month: 'gap-y-4 overflow-x-hidden w-full',
        month_caption: 'flex justify-center h-7 mx-10 relative items-center',
        month_grid: 'mt-4',
        months: 'flex flex-col relative',
        nav: 'flex items-start',
        outside:
          'day-outside text-muted-foreground opacity-50 aria-selected:bg-accent/50 aria-selected:text-muted-foreground aria-selected:opacity-30',
        range_end: 'day-range-end rounded-e-md',
        range_middle:
          'aria-selected:bg-accent hover:aria-selected:!bg-accent rounded-none aria-selected:text-accent-foreground hover:aria-selected:text-accent-foreground',
        range_start: 'day-range-start rounded-s-md',
        selected:
          'bg-primary text-primary-foreground hover:!bg-primary hover:text-primary-foreground focus:bg-primary focus:text-primary-foreground',
        today: 'bg-accent text-accent-foreground',
        week: 'flex w-full mt-2',
        weekday: 'text-muted-foreground w-8 font-normal text-[0.8rem]',
        weekdays: 'flex flex-row',
        ...classNames
      }}
      components={{
        Chevron: ({ orientation }) => {
          const Icon = orientation === 'left' ? ChevronLeftIcon : ChevronRightIcon
          return <Icon className='h-4 w-4' />
        }
      }}
      showOutsideDays={showOutsideDays}
      {...props}
    />
  )
}
Calendar.displayName = 'Calendar'

export { Calendar }
