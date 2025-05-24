'use client'

import { Link } from '@radix-ui/react-navigation-menu'
import { ChevronRight } from 'lucide-react'

import type { PageHeaderProps } from '#layouts/page-header/page-header'

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger
} from '#components/shadcn/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem
} from '#components/shadcn/sidebar'
import { cn } from '#lib/utils'

export function SidebarLinks({ pathname, siteRoutes }: PageHeaderProps) {
  const sections = Array.from(new Set(siteRoutes.map((route) => route.section)))

  return (
    <>
      {sections.map((section, i) => {
        return (
          <SidebarGroup key={i}>
            <SidebarGroupLabel>{section}</SidebarGroupLabel>
            <SidebarMenu>
              {siteRoutes
                .filter((rootRoute) => rootRoute.section === section)
                .map((rootRoute, i) => {
                  const isActive = rootRoute.children?.some((route) =>
                    pathname.startsWith(route.href)
                  )

                  if (!rootRoute.children?.length) {
                    return (
                      <SidebarMenuItem key={i}>
                        <Link href={rootRoute.href}>
                          <SidebarMenuButton
                            className={
                              pathname === rootRoute.href
                                ? 'bg-sidebar-accent'
                                : ''
                            }
                            tooltip={rootRoute.title}
                          >
                            <rootRoute.Icon
                              className={cn(
                                pathname == rootRoute.href
                                  ? 'text-sidebar-foreground'
                                  : ''
                              )}
                            />

                            <span>{rootRoute.title}</span>
                          </SidebarMenuButton>
                        </Link>
                      </SidebarMenuItem>
                    )
                  }

                  return (
                    <Collapsible
                      asChild
                      className='group/collapsible'
                      defaultOpen={isActive ?? false}
                      key={rootRoute.title}
                    >
                      <SidebarMenuItem>
                        <CollapsibleTrigger asChild>
                          <SidebarMenuButton tooltip={rootRoute.title}>
                            <rootRoute.Icon
                              className={cn(
                                pathname == rootRoute.href
                                  ? 'text-sidebar-foreground'
                                  : ''
                              )}
                            />
                            <span>{rootRoute.title}</span>
                            <ChevronRight className='ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
                          </SidebarMenuButton>
                        </CollapsibleTrigger>
                        <CollapsibleContent>
                          <SidebarMenuSub>
                            {rootRoute.children.map((route) => (
                              <SidebarMenuSubItem key={route.title}>
                                <SidebarMenuSubButton
                                  asChild
                                  className={
                                    pathname === route.href
                                      ? 'bg-sidebar-accent'
                                      : ''
                                  }
                                >
                                  <Link href={route.href}>
                                    <span>{route.title}</span>
                                  </Link>
                                </SidebarMenuSubButton>
                              </SidebarMenuSubItem>
                            ))}
                          </SidebarMenuSub>
                        </CollapsibleContent>
                      </SidebarMenuItem>
                    </Collapsible>
                  )
                })}
            </SidebarMenu>
          </SidebarGroup>
        )
      })}
    </>
  )
}
