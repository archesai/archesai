import { useLocation } from '@tanstack/react-router'

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
  siteRoutes: SiteRoute[]
}

export const PageHeader = ({ siteRoutes }: PageHeaderProps) => {
  const location = useLocation()

  // Split the pathname into segments and create breadcrumbs
  const pathSegments = location.pathname.split('/').filter(Boolean)

  const breadcrumbs = pathSegments.map((segment, index) => {
    const path = '/' + pathSegments.slice(0, index + 1).join('/')
    const title = segment
      .split('-')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')

    return {
      isLast: index === pathSegments.length - 1,
      path,
      title
    }
  })

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
              {breadcrumbs.map((breadcrumb, index) => (
                <div
                  className='flex items-center'
                  key={breadcrumb.path}
                >
                  {index > 0 && (
                    <BreadcrumbSeparator className='hidden md:block' />
                  )}
                  <BreadcrumbItem
                    className={index === 0 ? 'hidden md:block' : ''}
                  >
                    {breadcrumb.isLast ?
                      <BreadcrumbPage>{breadcrumb.title}</BreadcrumbPage>
                    : <BreadcrumbLink href={breadcrumb.path}>
                        {breadcrumb.title}
                      </BreadcrumbLink>
                    }
                  </BreadcrumbItem>
                </div>
              ))}
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
