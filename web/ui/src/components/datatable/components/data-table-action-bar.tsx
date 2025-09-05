'use no memo'

import type { JSX } from 'react'
import type { Table } from '@tanstack/react-table'

import { useCallback, useEffect, useLayoutEffect, useState } from 'react'
import { AnimatePresence, motion } from 'motion/react'
import * as ReactDOM from 'react-dom'

import { Loader2Icon, XCircleIcon } from '#components/custom/icons'
import { Button } from '#components/shadcn/button'
import { Separator } from '#components/shadcn/separator'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger
} from '#components/shadcn/tooltip'
import { cn } from '#lib/utils'

interface DataTableActionBarActionProps
  extends React.ComponentProps<typeof Button> {
  isPending?: boolean
  tooltip?: string
}

interface DataTableActionBarProps<TData>
  extends React.ComponentProps<typeof motion.div> {
  container?: DocumentFragment | Element | null
  table: Table<TData>
  visible?: boolean
}

interface DataTableActionBarSelectionProps<TData> {
  table: Table<TData>
}

function DataTableActionBar<TData>({
  children,
  className,
  container: containerProp,
  table,
  visible: visibleProp,
  ...props
}: DataTableActionBarProps<TData>): JSX.Element | null {
  const [mounted, setMounted] = useState(false)

  useLayoutEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        table.toggleAllRowsSelected(false)
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
    }
  }, [table])

  const container = containerProp ?? (mounted ? globalThis.document.body : null)

  if (!container) return null

  const visible =
    visibleProp ?? table.getFilteredSelectedRowModel().rows.length > 0

  return ReactDOM.createPortal(
    <AnimatePresence>
      {visible && (
        <motion.div
          animate={{ opacity: 1, y: 0 }}
          aria-orientation='horizontal'
          className={cn(
            'fixed inset-x-0 bottom-6 z-50 mx-auto flex w-fit flex-wrap items-center justify-center gap-2 rounded-md border bg-background p-2 text-foreground shadow-sm',
            className
          )}
          exit={{ opacity: 0, y: 20 }}
          initial={{ opacity: 0, y: 20 }}
          role='toolbar'
          transition={{ duration: 0.2, ease: 'easeInOut' }}
          {...props}
        >
          {children}
        </motion.div>
      )}
    </AnimatePresence>,
    container
  )
}

function DataTableActionBarAction({
  children,
  className,
  disabled,
  isPending,
  size = 'sm',
  tooltip,
  ...props
}: DataTableActionBarActionProps): JSX.Element {
  const trigger = (
    <Button
      className={cn(
        'gap-1.5 border border-secondary bg-secondary/50 hover:bg-secondary/70 [&>svg]:size-3.5',
        size === 'icon' ? 'size-7' : 'h-7',
        className
      )}
      disabled={disabled ?? isPending}
      size={size}
      variant='secondary'
      {...props}
    >
      {isPending ?
        <Loader2Icon className='animate-spin' />
      : children}
    </Button>
  )

  if (!tooltip) return trigger

  return (
    <Tooltip>
      <TooltipTrigger asChild>{trigger}</TooltipTrigger>
      <TooltipContent
        className='border bg-accent font-semibold text-foreground dark:bg-zinc-900 [&>span]:hidden'
        sideOffset={6}
      >
        <p>{tooltip}</p>
      </TooltipContent>
    </Tooltip>
  )
}

function DataTableActionBarSelection<TData>({
  table
}: DataTableActionBarSelectionProps<TData>): JSX.Element {
  const onClearSelection = useCallback(() => {
    table.toggleAllRowsSelected(false)
  }, [table])

  return (
    <div className='flex h-7 items-center rounded-md border pr-1 pl-2.5'>
      <span className='text-xs whitespace-nowrap'>
        {table.getFilteredSelectedRowModel().rows.length} selected
      </span>
      <Separator
        className='mr-1 ml-2 data-[orientation=vertical]:h-4'
        orientation='vertical'
      />
      <Tooltip>
        <TooltipTrigger asChild>
          <Button
            className='size-5'
            onClick={onClearSelection}
            size='icon'
            variant='ghost'
          >
            <XCircleIcon className='size-3.5' />
          </Button>
        </TooltipTrigger>
        <TooltipContent
          className='flex items-center gap-2 border bg-accent px-2 py-1 font-semibold text-foreground dark:bg-zinc-900 [&>span]:hidden'
          sideOffset={10}
        >
          <p>Clear selection</p>
          <kbd className='rounded border bg-background px-1.5 py-px font-mono text-[0.7rem] font-normal text-foreground shadow-xs select-none'>
            <abbr
              className='no-underline'
              title='Escape'
            >
              Esc
            </abbr>
          </kbd>
        </TooltipContent>
      </Tooltip>
    </div>
  )
}

export {
  DataTableActionBar,
  DataTableActionBarAction,
  DataTableActionBarSelection
}
