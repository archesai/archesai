import type { JSX } from 'react'

import type { SiteRoute } from '#lib/site-config.interface'

import { UserButton } from '#components/custom/user-button'
import { Separator } from '#components/shadcn/separator'
import { SidebarTrigger } from '#components/shadcn/sidebar'
import { BreadCrumbs } from '#layouts/page-header/components/breadcrumbs'
import { CommandMenu } from '#layouts/page-header/components/command-menu'
import { ThemeToggle } from '#layouts/page-header/components/theme-toggle'

export interface PageHeaderProps {
  siteRoutes: SiteRoute[]
}

export const PageHeader = ({ siteRoutes }: PageHeaderProps): JSX.Element => {
  return (
    <header className='flex h-16 shrink-0 justify-between px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12'>
      <div className='flex items-center justify-start gap-2'>
        <SidebarTrigger />
        <Separator
          className='data-[orientation=vertical]:h-4'
          orientation='vertical'
        />
        <BreadCrumbs />
      </div>
      <div className='flex items-center justify-end gap-2'>
        <CommandMenu siteRoutes={siteRoutes} />
        <ThemeToggle />
        <UserButton size={'sm'} />
      </div>
    </header>
  )
}
