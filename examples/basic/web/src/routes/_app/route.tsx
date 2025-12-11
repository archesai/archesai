import {
  AppSidebarContainer,
  PageHeaderContainer,
  SidebarInset,
  SidebarProvider,
} from "@archesai/ui";
import {
  createFileRoute,
  Outlet,
  redirect,
  useLocation,
  useNavigate,
} from "@tanstack/react-router";
import type { JSX } from "react";
import { getEnvConfig } from "#lib/config";
import { siteRoutes } from "#lib/site-config";

export const Route = createFileRoute("/_app")({
  beforeLoad: ({ context }) => {
    const { authEnabled } = getEnvConfig();
    if (authEnabled && !context.session?.data) {
      throw redirect({
        to: "/auth/login",
      });
    }
    return {
      session: context.session,
    };
  },
  component: AppLayout,
});

function AppLayout(): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  return (
    <SidebarProvider
      className="flex flex-col"
      defaultOpen={false}
    >
      <PageHeaderContainer
        location={location}
        navigate={navigate}
        siteRoutes={siteRoutes}
      />
      <div className="flex flex-1">
        <AppSidebarContainer
          location={location}
          navigate={navigate}
          siteRoutes={siteRoutes}
        />
        <SidebarInset className="flex flex-col">
          <main className="flex flex-1 flex-col p-4">
            <Outlet />
          </main>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
