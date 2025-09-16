import { Separator, SidebarInset, SidebarProvider } from "@archesai/ui";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import type { JSX } from "react";
import { AppSidebarContainer } from "#components/layouts/app-sidebar-container";
import { PageHeaderContainer } from "#components/layouts/page-header-container";
import { siteRoutes } from "#lib/site-config";

export const Route = createFileRoute("/_app")({
  beforeLoad: ({ context }) => {
    if (process.env.ARCHESAI_AUTH_ENABLED && !context.session?.data) {
      throw redirect({ to: "/auth/login" });
    }
  },
  component: AppLayout,
});

function AppLayout(): JSX.Element {
  return (
    <SidebarProvider defaultOpen={false}>
      {/* This is the sidebar that is displayed on the left side of the screen. */}
      <AppSidebarContainer siteRoutes={siteRoutes} />
      {/* This is the main content area. */}
      <SidebarInset className="max-h-screen">
        <PageHeaderContainer siteRoutes={siteRoutes} />
        <Separator />
        <div className="flex flex-1 flex-col overflow-y-auto p-4">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
