"use client";

import { NavMain } from "@/components/sidebar/nav-main";
import { NavUser } from "@/components/sidebar/nav-user";
import { RecentThreads } from "@/components/sidebar/recent-threads";
import { TeamSwitcher } from "@/components/team-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";

import { CreditQuota } from "../credit-quota";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="flex h-14 items-center justify-center">
        <TeamSwitcher />
      </SidebarHeader>
      <hr />
      <SidebarContent>
        <NavMain />
        <RecentThreads />
      </SidebarContent>
      <SidebarFooter>
        <CreditQuota />
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
