'use client'

import { useSidebar } from '@/components/ui/sidebar'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { siteConfig } from '@/config/site'
import { Menu } from 'lucide-react'
import { usePathname, useRouter } from 'next/navigation'

import { Button } from '../../ui/button'
import { CommandMenu } from './command-menu'
import { VerifyEmailAlert } from './email-verify'
import { ModeToggle } from './mode-toggle'
import { TitleAndDescription } from './title-and-description'
import { UserButton } from '../../user-button'
import { cn } from '@/lib/utils'

export const PageHeader = () => {
  const { toggleSidebar } = useSidebar()
  const router = useRouter()
  const pathname = usePathname() as string

  // combine all the routes from siteConfig
  const routes = siteConfig.routes
    .map((route) => [route, ...(route.children || [])])
    .flat()

  // find the current route
  const currentRoute = routes.find((route) => pathname === route.href)
  // get the title and description from the current route
  const title = currentRoute?.title
  const description = currentRoute?.description
  const Icon = currentRoute?.Icon

  const currentTabs = siteConfig.routes
    .find((route) => pathname.startsWith(route.href))
    ?.children?.filter((tab: any) => tab?.showInTabs)
  const activeTab = currentTabs?.find((tab) => pathname === tab.href)?.href

  return (
    <>
      <VerifyEmailAlert />

      <header className='flex w-full items-center justify-between bg-sidebar p-3 py-3'>
        <Button
          className='mr-3 flex h-8 w-8 text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-foreground'
          onClick={toggleSidebar}
          size='icon'
          variant='ghost'
        >
          <Menu className='h-5 w-5' />
        </Button>
        <div className='flex flex-1 items-center justify-end gap-3'>
          <CommandMenu />
          <ModeToggle />
          <div className='h-8 w-8'>
            <UserButton
              size={'sm'}
              side='bottom'
            />
          </div>
        </div>
      </header>

      {!currentTabs || currentTabs.length === 0 ? (
        <div className='border-b border-b-sidebar-accent' />
      ) : (
        <Tabs value={activeTab}>
          <TabsList className='h-8 w-full items-end justify-start rounded-none border-b bg-sidebar'>
            {currentTabs.map((tab) => {
              const isActive = tab.href === activeTab
              return (
                <TabsTrigger
                  className={cn(
                    `relative h-8 font-normal shadow-none transition-all hover:bg-sidebar-accent hover:text-sidebar-foreground data-[state=active]:text-sidebar-foreground [&::after]:absolute [&::after]:bottom-0 [&::after]:left-0 [&::after]:h-0.5 [&::after]:bg-blue-600 [&::after]:transition-all [&::after]:content-['']`,
                    `${isActive ? 'text-sidebar-foreground [&::after]:w-full' : 'text-muted-foreground [&::after]:w-0'}`
                  )}
                  key={tab.href}
                  onClick={() => {
                    router.push(tab.href)
                  }}
                  value={tab.href}
                >
                  {tab.title}
                </TabsTrigger>
              )
            })}
          </TabsList>
        </Tabs>
      )}

      <TitleAndDescription
        description={description}
        Icon={Icon}
        title={title}
      />
    </>
  )
}
