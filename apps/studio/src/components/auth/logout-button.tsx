import { Button, DropdownMenuItem, LogOutIcon } from "@archesai/ui";
import { useQueryClient } from "@tanstack/react-query";
import { useNavigate, useRouteContext } from "@tanstack/react-router";
import type { JSX, ReactNode } from "react";
import { useCallback } from "react";
import { useDeleteSession } from "#lib/index";

export interface LogoutButtonProps {
  variant?: "button" | "dropdown-item" | "icon-only";
  size?: "icon" | "default" | "sm" | "lg";
  className?: string;
  children?: ReactNode;
  showIcon?: boolean;
  onLogoutComplete?: () => void;
}

export function LogoutButton({
  variant = "button",
  size = "default",
  className,
  children,
  showIcon = true,
  onLogoutComplete,
}: LogoutButtonProps): JSX.Element {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { session } = useRouteContext({ from: "__root__" });
  const { mutateAsync: deleteSession } = useDeleteSession();

  const sessionData = session?.data;

  const handleLogout = useCallback(async () => {
    try {
      if (sessionData?.id) {
        await deleteSession({ id: sessionData.id });
      }
      queryClient.clear();
      await navigate({
        to: "/auth/login",
      });
      onLogoutComplete?.();
    } catch (error) {
      console.error("Logout error:", error);
      await navigate({
        to: "/auth/login",
      });
      onLogoutComplete?.();
    }
  }, [sessionData?.id, deleteSession, queryClient, navigate, onLogoutComplete]);

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

export function LogoutButtonWithShortcut(): JSX.Element {
  return (
    <LogoutButton variant="dropdown-item">
      <span className="flex flex-1 items-center justify-between">
        <span>Log out</span>
        <span className="ml-auto text-xs tracking-widest opacity-60">⌘Q</span>
      </span>
    </LogoutButton>
  );
}
