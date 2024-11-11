"use client";

import { NavMain } from "@/components/sidebar/nav-main";
import { NavUser } from "@/components/sidebar/nav-user";
// import { RecentLabels } from "@/components/sidebar/recent-labels";
import { OrganizationSwitcher } from "@/components/sidebar/organization-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";

import { CreditQuota } from "./credit-quota";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="flex h-14 items-center justify-center">
        <OrganizationSwitcher />
      </SidebarHeader>
      <hr />
      <SidebarContent className="gap-0">
        <NavMain />
        {/* <RecentLabels /> */}
      </SidebarContent>
      <SidebarFooter>
        <CreditQuota />
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
