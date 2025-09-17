import { SidebarInset, SidebarProvider } from "@archesai/ui";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import type { JSX } from "react";
import { AppSidebarContainer } from "#components/containers/app-sidebar-container";
import { PageHeaderContainer } from "#components/containers/page-header-container";
import { siteRoutes } from "#lib/site-config";

export const Route = createFileRoute("/_app")({
  beforeLoad: ({ context }) => {
    if (process.env.ARCHESAI_AUTH_ENABLED && !context.session?.data) {
      throw redirect({
        to: "/auth/login",
      });
    }
  },
  component: AppLayout,
});

function AppLayout(): JSX.Element {
  return (
    <SidebarProvider defaultOpen={true}>
      <AppSidebarContainer siteRoutes={siteRoutes} />
      <SidebarInset className="flex h-screen flex-col">
        <PageHeaderContainer siteRoutes={siteRoutes} />
        <main className="flex-1 overflow-auto p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
