import { Link, useLocation, useRouter } from "@tanstack/react-router";
import type { JSX } from "react";
import { ChevronRightIcon } from "#components/custom/icons";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "#components/shadcn/collapsible";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "#components/shadcn/sidebar";
import type { PageHeaderProps } from "#layouts/page-header/page-header";

export function SidebarLinks({ siteRoutes }: PageHeaderProps): JSX.Element {
  const router = useRouter();
  const pathname = useLocation().pathname;
  const sections = Array.from(
    new Set(siteRoutes.map((route) => route.section)),
  );
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
                  .map((rootRoute) => {
                    const isActive = rootRoute.children?.some((route) =>
                      router.state.location.pathname.startsWith(route.href),
                    );
                    const children = rootRoute.children ?? [];

                    if (!children.length) {
                      return (
                        <SidebarMenuItem key={rootRoute.href}>
                          <SidebarMenuButton
                            asChild
                            isActive={rootRoute.href === pathname}
                            tooltip={rootRoute.title}
                          >
                            <Link
                              className="text-muted-foreground"
                              to={rootRoute.href}
                            >
                              <rootRoute.Icon />
                              <span>{rootRoute.title}</span>
                            </Link>
                          </SidebarMenuButton>
                        </SidebarMenuItem>
                      );
                    }

                    return (
                      <Collapsible
                        asChild
                        className="group/collapsible"
                        defaultOpen={isActive ?? false}
                        key={rootRoute.title}
                      >
                        <SidebarMenuItem>
                          <CollapsibleTrigger asChild>
                            <SidebarMenuButton
                              className="text-muted-foreground"
                              tooltip={rootRoute.title}
                            >
                              <rootRoute.Icon />
                              <span>{rootRoute.title}</span>
                              <ChevronRightIcon className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
                            </SidebarMenuButton>
                          </CollapsibleTrigger>
                          <CollapsibleContent>
                            <SidebarMenuSub>
                              {children.map((route) => (
                                <SidebarMenuSubItem key={route.title}>
                                  <SidebarMenuSubButton
                                    asChild
                                    className="text-muted-foreground"
                                    isActive={route.href === pathname}
                                  >
                                    <Link to={route.href}>
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
            </SidebarGroupContent>
          </SidebarGroup>
        );
      })}
    </>
  );
}
