'use client'

import { usePathname } from 'next/navigation'

import { Authenticated } from '@archesai/ui/components/custom/authenticated'
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
  const pathname = usePathname()
  return (
    <>
      <Authenticated />
      <SidebarProvider>
        {/* This is the sidebar that is displayed on the left side of the screen. */}
        <AppSidebar
          pathname={pathname}
          siteRoutes={siteRoutes}
        />
        {/* This is the main content area. */}
        <SidebarInset>
          <main className='flex h-svh flex-col'>
            <PageHeader
              pathname={pathname}
              siteRoutes={siteRoutes}
            />
            <div className='container flex-1 overflow-auto p-4'>{children}</div>
          </main>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
