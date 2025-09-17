import { AppSidebar, OrganizationButton } from "@archesai/ui";
import type { SiteRoute } from "@archesai/ui/lib/site-config.interface";
import { useLocation, useNavigate } from "@tanstack/react-router";
import { UserButtonContainer } from "../navigation/user-button-container";

interface AppSidebarContainerProps {
  siteRoutes: SiteRoute[];
}

export function AppSidebarContainer({ siteRoutes }: AppSidebarContainerProps) {
  const location = useLocation();
  const navigate = useNavigate();

  const handleNavigate = (href: string) => {
    navigate({
      to: href,
    });
  };

  const handleSearch = (query: string) => {
    // Implement search logic here
    console.log("Search:", query);
  };

  return (
    <AppSidebar
      currentPath={location.pathname}
      onNavigate={handleNavigate}
      onSearch={handleSearch}
      organizationSlot={<OrganizationButton />}
      siteRoutes={siteRoutes}
      userMenuSlot={<UserButtonContainer />}
    />
  );
}
