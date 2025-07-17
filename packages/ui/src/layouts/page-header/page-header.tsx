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
      <header className='z-1 flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12'>
        <div className='flex w-full items-center gap-2 px-4'>
          <SidebarTrigger className='-ml-1 text-muted-foreground' />
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
    </>
  )
}
