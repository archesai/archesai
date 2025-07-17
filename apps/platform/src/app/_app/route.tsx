import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { createIsomorphicFn, createServerFn } from '@tanstack/react-start'
import { getHeaders as getServerHeaders } from '@tanstack/react-start/server'

import type { GetSession200 } from '@archesai/client'

import { getGetSessionQueryOptions } from '@archesai/client'
import {
  SidebarInset,
  SidebarProvider
} from '@archesai/ui/components/shadcn/sidebar'
import { AppSidebar } from '@archesai/ui/layouts/app-sidebar/app-sidebar'
import { PageHeader } from '@archesai/ui/layouts/page-header/page-header'

import { siteRoutes } from '#lib/site-config'

export const getHeaders = createServerFn({ method: 'GET' }).handler(() => {
  return getServerHeaders()
})

export const getHeadersIsomorphic = createIsomorphicFn()
  .client(() => {
    return getHeaders()
  })
  .server(() => {
    return getServerHeaders()
  })

export const Route = createFileRoute('/_app')({
  beforeLoad: async ({ context, location }) => {
    try {
      const headers = await getHeadersIsomorphic()
      const res = (await context.queryClient.fetchQuery(
        getGetSessionQueryOptions({
          request: {
            headers: [['cookie', headers.cookie ?? '']]
          }
        })
      )) as GetSession200 | null | undefined
      if (!res) {
        throw new Error('Session not found')
      }
    } catch (error) {
      console.error('Error fetching session:', error)
      redirect({
        search: location.search,
        throw: true,
        to: '/auth/login'
      })
    }
  },
  component: AppLayout
})

export default function AppLayout() {
  return (
    <SidebarProvider>
      {/* This is the sidebar that is displayed on the left side of the screen. */}
      <AppSidebar siteRoutes={siteRoutes} />
      {/* This is the main content area. */}
      <SidebarInset className='max-h-screen bg-primary/5 dark:bg-black/70'>
        <PageHeader siteRoutes={siteRoutes} />
        <div className='flex flex-1 flex-col gap-4 overflow-y-auto p-4'>
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
