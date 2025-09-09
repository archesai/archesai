import { Separator } from "@archesai/ui/components/shadcn/separator";
import {
  SidebarInset,
  SidebarProvider,
} from "@archesai/ui/components/shadcn/sidebar";
import { AppSidebar } from "@archesai/ui/layouts/app-sidebar/app-sidebar";
import { PageHeader } from "@archesai/ui/layouts/page-header/page-header";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import type { JSX } from "react";

import { siteRoutes } from "#lib/site-config";

export const Route = createFileRoute("/_app")({
  beforeLoad: ({ context }) => {
    if (!context.session?.data) {
      throw redirect({ to: "/auth/login" });
    }
  },
  component: AppLayout,
});

export default function AppLayout(): JSX.Element {
  return (
    <SidebarProvider>
      {/* This is the sidebar that is displayed on the left side of the screen. */}
      <AppSidebar siteRoutes={siteRoutes} />
      {/* This is the main content area. */}
      <SidebarInset className="max-h-screen">
        <PageHeader siteRoutes={siteRoutes} />
        <Separator />
        <div className="flex flex-1 flex-col overflow-y-auto p-4">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
