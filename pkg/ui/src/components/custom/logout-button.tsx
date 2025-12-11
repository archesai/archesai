import { Button, DropdownMenuItem, LogOutIcon } from "@archesai/ui";
import type { JSX, ReactNode } from "react";
import { useCallback } from "react";

export interface LogoutButtonProps {
  variant?: "button" | "dropdown-item" | "icon-only";
  size?: "icon" | "default" | "sm" | "lg";
  className?: string;
  children?: ReactNode;
  showIcon?: boolean;
  onLogout: () => void | Promise<void>;
}

export function LogoutButton({
  variant = "button",
  size = "default",
  className,
  children,
  showIcon = true,
  onLogout,
}: LogoutButtonProps): JSX.Element {
  const handleLogout = useCallback(async () => {
    try {
      await onLogout();
    } catch (error) {
      console.error("Logout error:", error);
    }
  }, [onLogout]);

  const iconElement = showIcon ? <LogOutIcon className="h-4 w-4" /> : null;
  const labelElement = children || "Log out";

  switch (variant) {
    case "dropdown-item":
      return (
        <DropdownMenuItem onClick={handleLogout}>
          {showIcon && <LogOutIcon className="mr-2 h-4 w-4" />}
          {labelElement}
        </DropdownMenuItem>
      );

    case "icon-only":
      return (
        <Button
          aria-label="Log out"
          className={className}
          onClick={handleLogout}
          size="icon"
          variant="ghost"
        >
          <LogOutIcon className="h-4 w-4" />
        </Button>
      );

    default:
      return (
        <Button
          className={className}
          onClick={handleLogout}
          size={size}
          variant="outline"
        >
          {iconElement && (
            <span className={children ? "mr-2" : ""}>{iconElement}</span>
          )}
          {labelElement}
        </Button>
      );
  }
}

export function LogoutButtonWithShortcut({
  onLogout,
}: {
  onLogout: () => void | Promise<void>;
}): JSX.Element {
  return (
    <LogoutButton
      onLogout={onLogout}
      variant="dropdown-item"
    >
      <span className="flex flex-1 items-center justify-between">
        <span>Log out</span>
        <span className="ml-auto text-xs tracking-widest opacity-60">⌘Q</span>
      </span>
    </LogoutButton>
  );
}
