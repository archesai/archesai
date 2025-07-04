import { useLocation } from '@tanstack/react-router'

import {
  SidebarInset,
  SidebarProvider
} from '@archesai/ui/components/shadcn/sidebar'
import { AppSidebar } from '@archesai/ui/layouts/app-sidebar/app-sidebar'
import { PageHeader } from '@archesai/ui/layouts/page-header/page-header'

import { siteRoutes } from '#lib/site-config'

export default function AppLayout({
  children
}: Readonly<{
  children: React.ReactNode
}>) {
  const location = useLocation()
  return (
    <>
      <SidebarProvider>
        {/* This is the sidebar that is displayed on the left side of the screen. */}
        <AppSidebar
          pathname={location.pathname}
          siteRoutes={siteRoutes}
        />
        {/* This is the main content area. */}
        <SidebarInset>
          <main className='flex h-svh flex-col'>
            <PageHeader
              pathname={location.pathname}
              siteRoutes={siteRoutes}
            />
            <div className='container flex-1 overflow-auto p-4'>{children}</div>
          </main>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
