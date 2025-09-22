import { SidebarInset, SidebarProvider } from "@archesai/ui";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import type { JSX } from "react";
import { AppSidebarContainer } from "#components/containers/app-sidebar-container";
import { PageHeaderContainer } from "#components/containers/page-header-container";
import { siteRoutes } from "#lib/site-config";

export const Route = createFileRoute("/_app")({
  beforeLoad: ({ context }) => {
    // if (!context.session?.data) {
    //   throw redirect({
    //     to: "/auth/login",
    //   });
    // } // FIXME
    return {
      session: context.session,
    };
  },
  component: AppLayout,
});

function AppLayout(): JSX.Element {
  return (
    <SidebarProvider defaultOpen={false}>
      <AppSidebarContainer siteRoutes={siteRoutes} />
      <SidebarInset>
        <PageHeaderContainer siteRoutes={siteRoutes} />
        <main className="flex-1 p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
