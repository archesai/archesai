// 'use client'

// import { Calendar as CalendarIcon } from 'lucide-react'

// import { Button } from '#components/shadcn/button'
// import { Calendar } from '#components/shadcn/calendar'
// import {
//   Popover,
//   PopoverContent,
//   PopoverTrigger
// } from '#components/shadcn/popover'
// import { cn } from '#lib/utils'

// import { useSearchQuery } from '#hooks/use-search-query'

export function DatePickerWithRange() {
  // const { searchQuery, setSearchQuery } = useSearchQuery()

  return (
    // <Popover>
    //   <PopoverTrigger asChild>
    //     <Button
    //       className={cn(
    //         'hidden h-8 justify-start gap-2 text-left font-normal md:flex',
    //         !range && 'text-muted-foreground'
    //       )}
    //       id='range'
    //       variant={'outline'}
    //     >
    //       <CalendarIcon className='h-5 w-5' />
    //       {range?.from ? (
    //         range.to ? (
    //           <>
    //             {format(range.from, 'LLL dd, y')} -{' '}
    //             {format(range.to, 'LLL dd, y')}
    //           </>
    //         ) : (
    //           format(range.from, 'LLL dd, y')
    //         )
    //       ) : (
    //         <span>Pick a date</span>
    //       )}
    //     </Button>
    //   </PopoverTrigger>
    //   <PopoverContent
    //     align='start'
    //     className='w-auto p-0'
    //   >
    //     <Calendar
    //       autoFocus
    //       defaultMonth={range?.from ?? new Date()}
    //       mode='range'
    //       numberOfMonths={2}
    //       onSelect={(range) => {
    //         if (range) {
    //           if (range.from && range.to) {
    //             setRange({
    //               from: range.from,
    //               to: range.to
    //             })
    //           }
    //         }
    //       }}
    //       selected={range}
    //     />
    //   </PopoverContent>
    // </Popover>
    <></>
  )
}
