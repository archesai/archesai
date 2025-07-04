import { Search } from 'lucide-react'

import type { PageHeaderProps } from '#layouts/page-header/page-header'

import { UserButton } from '#components/custom/user-button'
import { Label } from '#components/shadcn/label'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarInput,
  SidebarRail
} from '#components/shadcn/sidebar'
import { OrganizationButton } from '#layouts/app-sidebar/organization-button'
import { SidebarLinks } from '#layouts/app-sidebar/sidebar-links'

export function AppSidebar({
  pathname,
  siteRoutes,
  ...props
}: PageHeaderProps & React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <OrganizationButton />
        {/* <SearchForm /> */}
      </SidebarHeader>
      <SidebarContent>
        <SidebarLinks
          pathname={pathname}
          siteRoutes={siteRoutes}
        />
      </SidebarContent>
      <SidebarFooter>
        {/* <CreditQuota /> */}
        <UserButton />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

export function SearchForm({ ...props }: React.ComponentProps<'form'>) {
  return (
    <form {...props}>
      <SidebarGroup className='py-0'>
        <SidebarGroupContent className='relative'>
          <Label
            className='sr-only'
            htmlFor='search'
          >
            Search
          </Label>
          <SidebarInput
            className='pl-8'
            id='search'
            placeholder='Search the docs...'
          />
          <Search className='pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none' />
        </SidebarGroupContent>
      </SidebarGroup>
    </form>
  )
}
