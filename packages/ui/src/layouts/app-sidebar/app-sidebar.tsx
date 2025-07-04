import { useGetSession } from '@archesai/client'
import { UserEntity } from '@archesai/domain'

import type { PageHeaderProps } from '#layouts/page-header/page-header'

import { UserButton } from '#components/custom/user-button'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail
} from '#components/shadcn/sidebar'
import { OrganizationButton } from '#layouts/app-sidebar/organization-button'
import { SidebarLinks } from '#layouts/app-sidebar/sidebar-links'
import { CreditQuota } from './credit-usage'

export function AppSidebar({
  pathname,
  siteRoutes,
  ...props
}: PageHeaderProps & React.ComponentProps<typeof Sidebar>) {
  const { data: session, status } = useGetSession({
    fetch: {
      credentials: 'include'
    }
  })
  if (status !== 'success' || session.status === 401) {
    return null
  }

  return (
    <Sidebar
      collapsible='icon'
      {...props}
    >
      <SidebarHeader className='flex h-14 items-center justify-center'>
        <OrganizationButton user={new UserEntity(session.data as UserEntity)} />
      </SidebarHeader>
      <SidebarContent className='gap-0'>
        <SidebarLinks
          pathname={pathname}
          siteRoutes={siteRoutes}
        />
      </SidebarContent>
      <SidebarFooter>
        <CreditQuota />
        <UserButton />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
