import type { MouseEvent, ReactNode } from "react";
import { cn } from "./utils";

interface AnchorLinkProps {
  to: string;
  children: ReactNode;
  className?: string;
  onClick?: (e: MouseEvent<HTMLAnchorElement>) => void;
}

export const AnchorLink = ({
  to,
  children,
  className,
  onClick,
}: AnchorLinkProps) => {
  const handleClick = (event: MouseEvent<HTMLAnchorElement>) => {
    if (onClick) {
      onClick(event);
    }

    // Check if the hash is already in the URL
    const [pathname, hash] = to.split("#");
    if (
      hash &&
      pathname === window.location.pathname &&
      hash === window.location.hash.slice(1)
    ) {
      event.preventDefault();
      const target = document.getElementById(hash);
      if (target) {
        target.scrollIntoView({ behavior: "smooth" });
      }
    }
  };

  return (
    <a
      className={cn(className)}
      href={to}
      onClick={handleClick}
    >
      {children}
    </a>
  );
};
