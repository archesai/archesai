import type { AnchorHTMLAttributes, JSX } from "react";
import { cn } from "#lib/utils";

export interface LinkProps extends AnchorHTMLAttributes<HTMLAnchorElement> {
  href: string;
  children: React.ReactNode;
  isActive?: boolean;
  external?: boolean;
}

export function Link({
  href,
  children,
  className,
  isActive,
  external,
  onClick,
  ...props
}: LinkProps): JSX.Element {
  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    // If external link or has custom onClick, let it proceed normally
    if (external || onClick) {
      onClick?.(e);
      return;
    }

    // For internal links without onClick, prevent default to allow
    // platform app to handle navigation via event delegation if needed
    if (!external && !onClick) {
      // Let the platform handle this via a wrapper
    }
  };

  return (
    <a
      className={cn(className, isActive && "font-semibold")}
      href={href}
      onClick={handleClick}
      rel={external ? "noopener noreferrer" : undefined}
      target={external ? "_blank" : undefined}
      {...props}
    >
      {children}
    </a>
  );
}
