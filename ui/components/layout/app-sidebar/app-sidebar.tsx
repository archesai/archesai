import { NavUser } from "@/components/layout/app-sidebar/nav-user";
// import { RecentLabels } from "@/components/sidebar/recent-labels";
import { OrganizationSwitcher } from "@/components/layout/app-sidebar/organization-switcher";
import { SidebarLinks } from "@/components/layout/app-sidebar/sidebar-links";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";

import { CreditQuota } from "./credit-usage";

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader className="flex h-14 items-center justify-center">
        <OrganizationSwitcher />
      </SidebarHeader>
      <hr />
      <SidebarContent className="gap-0">
        <SidebarLinks />
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
