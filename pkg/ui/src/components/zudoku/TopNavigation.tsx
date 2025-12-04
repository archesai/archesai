import type { ReactNode } from "react";

export const TopNavigation = ({ children }: { children?: ReactNode }) => {
  return (
    <nav className="flex items-center gap-6 px-4 py-2 lg:px-8">{children}</nav>
  );
};
