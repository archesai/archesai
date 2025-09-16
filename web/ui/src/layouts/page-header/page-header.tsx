import type { JSX, ReactNode } from "react";
import { Separator } from "#components/shadcn/separator";
import { SidebarTrigger } from "#components/shadcn/sidebar";

export interface PageHeaderProps {
  breadcrumbs?: ReactNode;
  commandMenu?: ReactNode;
  themeToggle?: ReactNode;
  userMenu?: ReactNode;
  showSidebarTrigger?: boolean;
  className?: string;
}

export const PageHeader = ({
  breadcrumbs,
  commandMenu,
  themeToggle,
  userMenu,
  showSidebarTrigger = true,
  className = "flex h-16 shrink-0 justify-between px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12",
}: PageHeaderProps): JSX.Element => {
  return (
    <header className={className}>
      <div className="flex items-center justify-start gap-2">
        {showSidebarTrigger && (
          <>
            <SidebarTrigger />
            <Separator
              className="data-[orientation=vertical]:h-4"
              orientation="vertical"
            />
          </>
        )}
        {breadcrumbs}
      </div>
      <div className="flex items-center justify-end gap-2">
        {commandMenu}
        {themeToggle}
        {userMenu}
      </div>
    </header>
  );
};
