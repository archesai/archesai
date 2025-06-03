'use client'

import { Menu } from 'lucide-react'

import type { SiteRoute } from '#lib/site-config.interface'

import { UserButton } from '#components/custom/user-button'
import { Button } from '#components/shadcn/button'
import { useSidebar } from '#components/shadcn/sidebar'
import { Tabs, TabsList, TabsTrigger } from '#components/shadcn/tabs'
import { CommandMenu } from '#layouts/page-header/command-menu'
import { VerifyEmailAlert } from '#layouts/page-header/email-verify'
import { ModeToggle } from '#layouts/page-header/mode-toggle'
import { TitleAndDescription } from '#layouts/page-header/title-and-description'
import { cn } from '#lib/utils'

export interface PageHeaderProps {
  pathname: string
  siteRoutes: SiteRoute[]
}

export const PageHeader = ({ pathname, siteRoutes }: PageHeaderProps) => {
  const { toggleSidebar } = useSidebar()

  // find the current route
  const currentRoute = siteRoutes
    .map((route) => [route, ...(route.children ?? [])])
    .flat()
    .find((route) => pathname === route.href)
  if (!currentRoute) {
    return null
  }

  // Tabination
  const tabs = siteRoutes
    .find((route) => pathname.startsWith(route.href))
    ?.children?.filter((tab) => tab.showInTabs)
  if (!tabs || tabs.length === 0) {
    return <div className='border-b' />
  }
  const activeTab = tabs.find((tab) => pathname === tab.href)?.href
  if (!activeTab) {
    return null
  }

  return (
    <>
      <VerifyEmailAlert />

      <header className='flex w-full items-center justify-between bg-sidebar p-3 py-3'>
        <Button
          className='mr-3 flex h-8 w-8 border-sidebar-accent bg-sidebar text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-foreground'
          onClick={toggleSidebar}
          size='icon'
          variant='outline'
        >
          <Menu className='h-5 w-5' />
        </Button>
        <div className='flex flex-1 items-center justify-end gap-3'>
          <CommandMenu siteRoutes={siteRoutes} />
          <ModeToggle />
          <div className='h-8 w-8'>
            <UserButton
              side='bottom'
              size={'sm'}
            />
          </div>
        </div>
      </header>

      <Tabs value={activeTab}>
        <TabsList className='h-8 w-full items-end justify-start rounded-none border-b bg-sidebar'>
          {tabs.map((tab) => {
            const isActive = tab.href === activeTab
            return (
              <TabsTrigger
                className={cn(
                  `relative h-8 font-normal shadow-none transition-all hover:bg-sidebar-accent hover:text-sidebar-foreground data-[state=active]:bg-sidebar data-[state=active]:text-sidebar-foreground [&::after]:absolute [&::after]:bottom-0 [&::after]:left-0 [&::after]:h-0.5 [&::after]:bg-blue-600 [&::after]:transition-all [&::after]:content-['']`,
                  isActive ?
                    'text-sidebar-foreground [&::after]:w-full'
                  : 'text-muted-foreground [&::after]:w-0'
                )}
                key={tab.href}
                onClick={() => {
                  window.location.href = tab.href
                }}
                value={tab.href}
              >
                {tab.title}
              </TabsTrigger>
            )
          })}
        </TabsList>
      </Tabs>

      <TitleAndDescription siteRoute={currentRoute} />
    </>
  )
}
