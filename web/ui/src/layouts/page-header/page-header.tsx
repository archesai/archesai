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
}: PageHeaderProps): JSX.Element => {
  return (
    <header
      className={"sticky top-0 z-50 flex h-14 justify-between border-b px-2"}
    >
      <div className="flex items-center gap-2">
        {showSidebarTrigger && <SidebarTrigger />}
        <Separator orientation="vertical" />
        {breadcrumbs}
      </div>
      <div className="flex items-center gap-2 px-2">
        {commandMenu}
        {themeToggle}
        {userMenu}
      </div>
    </header>
  );
};
