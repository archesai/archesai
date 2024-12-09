import { Authenticated } from '@/components/authenticated'
import { AppSidebar } from '@/components/layout/app-sidebar/app-sidebar'
import { PageHeader } from '@/components/layout/page-header/page-header'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'

export default function AppLayout({
  children
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <>
      <Authenticated />
      <SidebarProvider>
        {/* This is the sidebar that is displayed on the left side of the screen. */}
        <AppSidebar />
        {/* This is the main content area. */}
        <SidebarInset>
          <main className='flex h-svh flex-col'>
            <PageHeader />
            <div className='flex-1 overflow-auto p-4'>{children}</div>
          </main>
        </SidebarInset>
      </SidebarProvider>
    </>
  )
}
