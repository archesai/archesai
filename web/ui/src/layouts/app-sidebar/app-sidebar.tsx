import type { JSX, ReactNode } from "react";
import { SearchIcon } from "#components/custom/icons";
import { Label } from "#components/shadcn/label";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInput,
} from "#components/shadcn/sidebar";
import { SidebarLinks } from "#layouts/app-sidebar/sidebar-links";
import type { SiteRoute } from "#lib/site-config.interface";

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  siteRoutes: SiteRoute[];
  currentPath?: string;
  onNavigate?: (href: string) => void;
  organizationSlot?: ReactNode;
  userMenuSlot?: ReactNode;
  onSearch?: (query: string) => void;
}

export function AppSidebar({
  siteRoutes,
  currentPath,
  onNavigate,
  organizationSlot,
  userMenuSlot,
  onSearch,
  ...props
}: AppSidebarProps): JSX.Element {
  return (
    <Sidebar
      {...props}
      collapsible="icon"
      variant="inset"
    >
      <SidebarHeader>
        {organizationSlot}
        {onSearch && <SearchForm onSubmit={onSearch} />}
      </SidebarHeader>
      <SidebarContent>
        <SidebarLinks
          currentPath={currentPath}
          onNavigate={onNavigate}
          siteRoutes={siteRoutes}
        />
      </SidebarContent>
      {userMenuSlot && <SidebarFooter>{userMenuSlot}</SidebarFooter>}
    </Sidebar>
  );
}

interface SearchFormProps
  extends Omit<React.ComponentProps<"form">, "onSubmit"> {
  onSubmit?: (query: string) => void;
}

export function SearchForm({
  onSubmit,
  ...props
}: SearchFormProps): JSX.Element {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const query = formData.get("search") as string;
    onSubmit?.(query);
  };

  return (
    <form
      {...props}
      className="p-2"
      onSubmit={handleSubmit}
    >
      <div className="relative">
        <Label
          className="sr-only"
          htmlFor="search"
        >
          Search
        </Label>
        <SearchIcon className="-translate-y-1/2 pointer-events-none absolute top-1/2 left-2 size-4 text-muted-foreground" />
        <SidebarInput
          className="pl-8"
          id="search"
          name="search"
          placeholder="Search the docs..."
        />
      </div>
    </form>
  );
}
