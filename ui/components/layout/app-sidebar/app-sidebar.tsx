import { UserButton } from '@/components/user-button'
import { OrganizationButton } from '@/components/layout/app-sidebar/organization-button'
import { SidebarLinks } from '@/components/layout/app-sidebar/sidebar-links'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail
} from '@/components/ui/sidebar'

import { CreditQuota } from './credit-usage'

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar
      collapsible='icon'
      {...props}
    >
      <SidebarHeader className='flex h-14 items-center justify-center'>
        <OrganizationButton />
      </SidebarHeader>
      <SidebarContent className='gap-0'>
        <SidebarLinks />
      </SidebarContent>
      <SidebarFooter>
        <CreditQuota />
        <UserButton />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
