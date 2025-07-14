import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

import { getGetSessionSuspenseQueryOptions } from '@archesai/client'
import {
  SidebarInset,
  SidebarProvider
} from '@archesai/ui/components/shadcn/sidebar'
import { AppSidebar } from '@archesai/ui/layouts/app-sidebar/app-sidebar'
import { PageHeader } from '@archesai/ui/layouts/page-header/page-header'

import { siteRoutes } from '#lib/site-config'

export const Route = createFileRoute('/_app')({
  component: AppLayout,
  loader: async ({ context, location }) => {
    try {
      await context.queryClient.ensureQueryData(
        getGetSessionSuspenseQueryOptions()
      )
    } catch (error) {
      console.error('Error fetching session:', error)
      // eslint-disable-next-line @typescript-eslint/only-throw-error
      throw redirect({
        search: location.search,
        to: '/auth/login'
      })
    }
  }
})

export default function AppLayout() {
  return (
    <SidebarProvider>
      {/* This is the sidebar that is displayed on the left side of the screen. */}
      <AppSidebar siteRoutes={siteRoutes} />
      {/* This is the main content area. */}
      <SidebarInset>
        <PageHeader siteRoutes={siteRoutes} />
        <div className='flex flex-1 flex-col gap-4 p-4'>
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
