'use client'

import type { DialogProps } from '@radix-ui/react-dialog'

import { useCallback, useEffect, useState } from 'react'
import { Laptop, Moon, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'

import type { SiteRoute } from '#lib/site-config.interface'

import { Button } from '#components/shadcn/button'
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator
} from '#components/shadcn/command'
import { DialogDescription, DialogTitle } from '#components/shadcn/dialog'
import { cn } from '#lib/utils'

export function CommandMenu({
  siteRoutes,
  ...props
}: DialogProps & {
  siteRoutes: SiteRoute[]
}) {
  const [open, setOpen] = useState(false)
  const { setTheme } = useTheme()

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.key === 'k' && (e.metaKey || e.ctrlKey)) || e.key === '/') {
        if (
          (e.target instanceof HTMLElement && e.target.isContentEditable) ||
          e.target instanceof HTMLInputElement ||
          e.target instanceof HTMLTextAreaElement ||
          e.target instanceof HTMLSelectElement
        ) {
          return
        }

        e.preventDefault()
        setOpen((open) => !open)
      }
    }

    document.addEventListener('keydown', down)
    return () => {
      document.removeEventListener('keydown', down)
    }
  }, [])

  const runCommand = useCallback((command: () => unknown) => {
    setOpen(false)
    command()
  }, [])

  return (
    <div>
      <Button
        className={cn(
          'h-8 w-full justify-between gap-2 rounded-lg border-sidebar-border bg-sidebar text-base font-normal text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground md:text-sm'
        )}
        onClick={() => {
          setOpen(true)
        }}
        variant='outline'
        {...props}
      >
        <span className='hidden sm:inline-flex'>
          Type a command or search...
        </span>
        <span className='inline-flex sm:hidden'>Search...</span>

        <kbd className='pointer-events-none flex h-5 items-center gap-1 rounded-xs border border-sidebar-accent bg-sidebar-accent p-2 font-mono text-[10px] font-medium select-none'>
          <span className='text-xs'>⌘</span>
          <span>K</span>
        </kbd>
      </Button>
      <CommandDialog
        onOpenChange={setOpen}
        open={open}
      >
        <div className='hidden'>
          <DialogTitle />
        </div>
        <DialogDescription />
        <CommandInput placeholder='Type a command or search...' />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          {siteRoutes.map((rootRoute) => (
            <CommandGroup
              heading={rootRoute.title}
              key={rootRoute.title}
            >
              {rootRoute.children?.map((route) => (
                <CommandItem
                  className='flex gap-2'
                  key={route.href}
                  onClick={() => {
                    runCommand(() => {
                      window.location.href = route.href
                    })
                  }}
                  onSelect={() => {
                    runCommand(() => {
                      window.location.href = route.href
                    })
                  }}
                  value={route.title}
                >
                  <route.Icon className='h-5 w-5' />
                  <span>{route.title}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          ))}
          <CommandSeparator />
          <CommandGroup heading='Theme'>
            {['light', 'dark', 'system'].map((theme) => (
              <CommandItem
                className='flex gap-2'
                key={theme}
                onSelect={() => {
                  runCommand(() => {
                    setTheme(theme)
                  })
                }}
              >
                {theme === 'light' && <Sun className='h-5 w-5' />}
                {theme === 'dark' && <Moon className='h-5 w-5' />}
                {theme === 'system' && <Laptop className='h-5 w-5' />}
                <span>{theme.charAt(0).toUpperCase() + theme.slice(1)}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </div>
  )
}
