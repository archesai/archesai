import type { JSX } from "react";
import { ChevronRightIcon } from "#components/custom/icons";
import { Link } from "#components/primitives/link";
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
import type { SiteRoute } from "#lib/site-config.interface";

export interface SidebarLinksProps {
  siteRoutes: SiteRoute[];
  currentPath?: string | undefined;
  onNavigate?: ((href: string) => void) | undefined;
}

export function SidebarLinks({
  siteRoutes,
  currentPath = "",
  onNavigate,
}: SidebarLinksProps): JSX.Element {
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
                      currentPath.startsWith(route.href),
                    );
                    const children = rootRoute.children ?? [];

                    if (!children.length) {
                      return (
                        <SidebarMenuItem
                          className="relative group-data-[collapsible=icon]:my-1"
                          key={rootRoute.href}
                        >
                          {rootRoute.href === currentPath && (
                            <div className="-ml-2 absolute top-0 left-0 h-full w-0.5 bg-primary group-data-[collapsible=icon]:hidden" />
                          )}
                          <SidebarMenuButton
                            asChild
                            isActive={rootRoute.href === currentPath}
                            tooltip={rootRoute.title}
                          >
                            <Link
                              className="text-muted-foreground"
                              href={rootRoute.href}
                              onClick={(e) => {
                                e.preventDefault();
                                onNavigate?.(rootRoute.href);
                              }}
                            >
                              <rootRoute.Icon />
                              {rootRoute.title}
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
                        <SidebarMenuItem className="group-data-[collapsible=icon]:my-1">
                          <CollapsibleTrigger asChild>
                            <SidebarMenuButton
                              className="flex items-center gap-2 text-muted-foreground group-data-[collapsible=icon]:gap-0"
                              tooltip={rootRoute.title}
                            >
                              <rootRoute.Icon className="group-data-[collapsible=icon]:mx-auto" />
                              <span className="group-data-[collapsible=icon]:hidden">
                                {rootRoute.title}
                              </span>
                              <ChevronRightIcon className="ml-auto transition-transform duration-200 group-data-[collapsible=icon]:hidden group-data-[state=open]/collapsible:rotate-90" />
                            </SidebarMenuButton>
                          </CollapsibleTrigger>
                          <CollapsibleContent>
                            <SidebarMenuSub>
                              {children.map((route) => (
                                <SidebarMenuSubItem key={route.title}>
                                  <SidebarMenuSubButton
                                    asChild
                                    className="text-muted-foreground"
                                    isActive={route.href === currentPath}
                                  >
                                    <Link
                                      href={route.href}
                                      onClick={(e) => {
                                        e.preventDefault();
                                        onNavigate?.(route.href);
                                      }}
                                    >
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
