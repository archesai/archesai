"use client";

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { siteConfig } from "@/config/site";
import { ChevronRight } from "lucide-react";
import Link from "next/link";

export function NavMain() {
  return (
    <SidebarGroup>
      <SidebarGroupLabel>Platform</SidebarGroupLabel>
      <SidebarMenu>
        {siteConfig.routes.map((rootRoute) => {
          const isActive = rootRoute.children.some((route) =>
            window.location.pathname.startsWith(route.href)
          );
          return (
            <Collapsible
              asChild
              className="group/collapsible"
              defaultOpen={isActive}
              key={rootRoute.title}
            >
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuButton tooltip={rootRoute.title}>
                    {rootRoute.Icon && <rootRoute.Icon />}
                    <span>{rootRoute.title}</span>
                    <ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    {rootRoute.children?.map((route) => (
                      <SidebarMenuSubItem key={route.title}>
                        <SidebarMenuSubButton
                          asChild
                          className={`${
                            window.location.pathname === route.href
                              ? "bg-muted"
                              : ""
                          }`}
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
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}
