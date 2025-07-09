import { ChevronRight } from 'lucide-react'

import type { PageHeaderProps } from '#layouts/page-header/page-header'

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger
} from '#components/shadcn/collapsible'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem
} from '#components/shadcn/sidebar'
import { useLinkComponent } from '#hooks/use-link'

export function SidebarLinks({ pathname, siteRoutes }: PageHeaderProps) {
  const sections = Array.from(new Set(siteRoutes.map((route) => route.section)))
  const Link = useLinkComponent()
  return (
    <>
      {sections.map((section) => {
        return (
          <SidebarGroup key={section}>
            <SidebarGroupLabel>{section}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {siteRoutes
                  .filter((rootRoute) => rootRoute.section === section)
                  .map((rootRoute, i) => {
                    const isActive = rootRoute.children?.some((route) =>
                      pathname.startsWith(route.href)
                    )
                    const children = rootRoute.children ?? []

                    if (!children.length) {
                      return (
                        <SidebarMenuItem key={i}>
                          <SidebarMenuButton
                            // isActive={item.isActive}
                            asChild
                            tooltip={rootRoute.title}
                          >
                            <Link href={rootRoute.href}>
                              <rootRoute.Icon />
                              <span>{rootRoute.title}</span>
                            </Link>
                          </SidebarMenuButton>
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
                              <rootRoute.Icon />
                              <span>{rootRoute.title}</span>
                              <ChevronRight className='ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90' />
                            </SidebarMenuButton>
                          </CollapsibleTrigger>
                          <CollapsibleContent>
                            <SidebarMenuSub>
                              {children.map((route) => (
                                <SidebarMenuSubItem key={route.title}>
                                  <SidebarMenuSubButton asChild>
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
            </SidebarGroupContent>
          </SidebarGroup>
        )
      })}
    </>
  )
}
