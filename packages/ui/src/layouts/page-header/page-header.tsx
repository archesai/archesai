import type { SiteRoute } from '#lib/site-config.interface'

import { UserButton } from '#components/custom/user-button'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '#components/shadcn/breadcrumb'
import { Separator } from '#components/shadcn/separator'
import { SidebarTrigger } from '#components/shadcn/sidebar'
import { CommandMenu } from '#layouts/page-header/command-menu'
import { ModeToggle } from '#layouts/page-header/mode-toggle'

export interface PageHeaderProps {
  pathname: string
  siteRoutes: SiteRoute[]
}

export const PageHeader = ({ siteRoutes }: PageHeaderProps) => {
  return (
    <>
      {/* <VerifyEmailAlert /> */}

      <header className='flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12'>
        <div className='flex w-full items-center gap-2 px-4'>
          <SidebarTrigger className='-ml-1' />
          <Separator
            className='mr-2 data-[orientation=vertical]:h-4'
            orientation='vertical'
          />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem className='hidden md:block'>
                <BreadcrumbLink href='#'>
                  Building Your Application
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator className='hidden md:block' />
              <BreadcrumbItem>
                <BreadcrumbPage>Data Fetching</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <div className='flex flex-1 items-center justify-end gap-3'>
            <CommandMenu siteRoutes={siteRoutes} />
            <ModeToggle />
            <div>
              <UserButton size={'sm'} />
            </div>
          </div>
        </div>
      </header>

      {/* {tabs[0]?.href && (
        <Tabs value={tabs[0].href}>
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
      )}

      <TitleAndDescription siteRoute={currentRoute} /> */}
    </>
  )
}
