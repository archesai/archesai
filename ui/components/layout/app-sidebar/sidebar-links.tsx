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
import { usePathname } from "next/navigation";

export function SidebarLinks() {
  const pathname = usePathname();
  const sections = new Set(siteConfig.routes.map((route) => route.section))
    .values()
    .toArray();

  return (
    <>
      {sections.map((section, i) => {
        return (
          <SidebarGroup key={i}>
            <SidebarGroupLabel>{section}</SidebarGroupLabel>
            <SidebarMenu>
              {siteConfig.routes
                .filter((rootRoute) => rootRoute.section === section)
                .map((rootRoute, i) => {
                  const isActive = rootRoute?.children?.some((route) =>
                    pathname.startsWith(route.href)
                  );

                  if (!rootRoute?.children?.length) {
                    return (
                      <SidebarMenuItem key={i}>
                        <Link href={rootRoute.href}>
                          <SidebarMenuButton
                            className={`${pathname === rootRoute.href ? "bg-muted" : ""}`}
                            tooltip={rootRoute.title}
                          >
                            {rootRoute.Icon && <rootRoute.Icon />}

                            <span>{rootRoute.title}</span>
                          </SidebarMenuButton>
                        </Link>
                      </SidebarMenuItem>
                    );
                  }

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
                                    pathname === route.href ? "bg-muted" : ""
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
      })}
    </>
  );
}
