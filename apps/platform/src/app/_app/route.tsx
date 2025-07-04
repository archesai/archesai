import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation
} from '@tanstack/react-router'

import { getSession } from '@archesai/client'
import {
  SidebarInset,
  SidebarProvider
} from '@archesai/ui/components/shadcn/sidebar'
import { AppSidebar } from '@archesai/ui/layouts/app-sidebar/app-sidebar'
import { PageHeader } from '@archesai/ui/layouts/page-header/page-header'

import { siteRoutes } from '#lib/site-config'

export const Route = createFileRoute('/_app')({
  beforeLoad: async ({ location }) => {
    try {
      const user = await getSession({
        credentials: 'include'
      })
      return user
    } catch {
      return redirect({
        search: {
          // Use the current location to power a redirect after login
          // (Do not use `router.state.resolvedLocation` as it can
          // potentially lag behind the actual current location)
          redirect: location.href
        },
        throw: true,
        to: '/auth/login'
      })
    }
  },
  component: AppLayout
})

export default function AppLayout() {
  const location = useLocation()
  return (
    <SidebarProvider>
      {/* This is the sidebar that is displayed on the left side of the screen. */}
      <AppSidebar
        pathname={location.pathname}
        siteRoutes={siteRoutes}
      />
      {/* This is the main content area. */}
      <SidebarInset>
        <PageHeader
          pathname={location.pathname}
          siteRoutes={siteRoutes}
        />
        <div className='flex flex-1 flex-col gap-4 p-4'>
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
