import type { ReactNode } from "react";

interface LinkProps {
  to?: string;
  href?: string;
  className?: string;
  children: ReactNode;
  onClick?: () => void;
}

export const Link = ({ to, href, className, children, onClick }: LinkProps) => {
  const url = to || href || "#";

  return (
    <a
      className={className}
      href={url}
      onClick={onClick}
    >
      {children}
    </a>
  );
};

export const NavLink = Link;
