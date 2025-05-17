import { getOneUser } from '@archesai/client'
import { UserEntity } from '@archesai/domain'

import type { SiteRoute } from '#lib/site-config.interface'

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

export async function AppSidebar({
  siteRoutes = [],
  ...props
}: React.ComponentProps<typeof Sidebar> & {
  siteRoutes?: SiteRoute[]
}) {
  const pathname = window.location.pathname
  const { data: user, status } = await getOneUser('me')
  if (status !== 200) {
    return null
  }
  return (
    <Sidebar
      collapsible='icon'
      {...props}
    >
      <SidebarHeader className='flex h-14 items-center justify-center'>
        <OrganizationButton
          user={
            new UserEntity({
              ...user.data,
              ...user.data.attributes
            })
          }
        />
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
