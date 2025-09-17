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
  className,
}: PageHeaderProps): JSX.Element => {
  const headerClassName =
    className ||
    [
      "flex h-14 shrink-0 items-center justify-between",
      "border-b px-4",
      "transition-[width,height] ease-linear",
      "group-has-data-[collapsible=icon]/sidebar-wrapper:h-12",
    ].join(" ");

  return (
    <header className={headerClassName}>
      <div className="flex items-center gap-2">
        {showSidebarTrigger && (
          <>
            <SidebarTrigger className="-ml-1" />
            <Separator
              className="h-4"
              orientation="vertical"
            />
          </>
        )}
        {breadcrumbs}
      </div>
      <div className="flex items-center gap-2">
        {commandMenu}
        {themeToggle}
        {userMenu}
      </div>
    </header>
  );
};
