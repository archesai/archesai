import type { JSX, ReactNode } from "react";
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
}: PageHeaderProps): JSX.Element => {
  return (
    <header className={"flex h-14 justify-between border-b px-2"}>
      <div className="flex items-center gap-2">
        {showSidebarTrigger && <SidebarTrigger />}
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
