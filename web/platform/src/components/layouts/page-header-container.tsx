import {
  BreadCrumbs,
  CommandMenu,
  PageHeader,
  ThemeToggle,
} from "@archesai/ui";
import type { SiteRoute } from "@archesai/ui/lib/site-config.interface";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { UserButtonContainer } from "../navigation/user-button-container";

interface PageHeaderContainerProps {
  siteRoutes: SiteRoute[];
}

export function PageHeaderContainer({ siteRoutes }: PageHeaderContainerProps) {
  const location = useLocation();
  const navigate = useNavigate();

  const handleNavigate = (path: string) => {
    navigate({ to: path });
  };

  return (
    <PageHeader
      breadcrumbs={
        <BreadCrumbs
          currentPath={location.pathname}
          onNavigate={handleNavigate}
        />
      }
      commandMenu={<CommandMenu siteRoutes={siteRoutes} />}
      themeToggle={<ThemeToggle />}
      userMenu={<UserButtonContainer />}
    />
  );
}
