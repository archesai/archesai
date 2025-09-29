import type { JSX } from "react";
import {
  BadgeCheckIcon,
  ChevronsUpDownIcon,
  CreditCardIcon,
  LogOutIcon,
  SparklesIcon,
} from "#components/custom/icons";
import { Avatar, AvatarFallback, AvatarImage } from "#components/shadcn/avatar";
import { Badge } from "#components/shadcn/badge";
import { Button } from "#components/shadcn/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "#components/shadcn/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "#components/shadcn/sidebar";
import { useIsMobile } from "#hooks/use-mobile";

export interface PureUserButtonProps {
  // User data
  user?: {
    email: string;
    id: string;
    image?: null | string;
    name: string;
  };

  // Organization data
  currentOrganizationId?: string | null;
  organizations?: Array<{
    id: string;
    name: string;
  }>;

  // Callbacks
  onLogout?: () => void;
  onNavigateToProfile?: () => void;
  onNavigateToBilling?: () => void;
  onSwitchOrganization?: (organizationId: string) => void;

  // UI props
  side?: "bottom" | "left" | "right" | "top";
  size?: "default" | "lg" | "sm" | null | undefined;
}

/**
 * Pure presentational UserButton component.
 * All business logic is handled by the container.
 */
export function PureUserButton({
  user,
  currentOrganizationId,
  organizations = [],
  onLogout,
  onNavigateToProfile,
  onNavigateToBilling,
  onSwitchOrganization,
  size = "lg",
}: PureUserButtonProps): JSX.Element {
  const isMobile = useIsMobile();
  const defaultOrgname = "Arches Platform";

  // If no user data is provided, render a placeholder
  if (!user) {
    return (
      <SidebarMenu>
        <SidebarMenuItem>
          <Button
            className="h-8 w-8"
            size="sm"
            variant="ghost"
          >
            <Avatar>
              <AvatarFallback>U</AvatarFallback>
            </Avatar>
          </Button>
        </SidebarMenuItem>
      </SidebarMenu>
    );
  }

  const userInitials = user.name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase();

  const currentOrg = organizations.find(
    (org) => org.id === currentOrganizationId,
  );

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size={size}
              tooltip={"Select Organization"}
            >
              <Avatar className="h-8 w-8 rounded-lg">
                <AvatarImage
                  alt={user.name}
                  src={user.image ?? undefined}
                />
                <AvatarFallback className="rounded-lg">
                  {userInitials}
                </AvatarFallback>
              </Avatar>
              <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                <span className="truncate font-semibold">{user.name}</span>
                <span className="truncate text-xs">
                  {currentOrg?.name ?? defaultOrgname}
                </span>
              </div>
              <ChevronsUpDownIcon className="ml-auto size-4 group-data-[collapsible=icon]:hidden" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="p-0 font-normal">
              <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                <Avatar className="h-8 w-8 rounded-lg">
                  <AvatarImage
                    alt={user.name}
                    src={user.image ?? undefined}
                  />
                  <AvatarFallback className="rounded-lg">
                    {userInitials}
                  </AvatarFallback>
                </Avatar>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">{user.name}</span>
                  <span className="truncate text-xs">{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              {organizations.length > 0 && (
                <>
                  <DropdownMenuSub>
                    <DropdownMenuSubTrigger>
                      <BadgeCheckIcon className="mr-2 h-4 w-4" />
                      Switch Organization
                    </DropdownMenuSubTrigger>
                    <DropdownMenuPortal>
                      <DropdownMenuSubContent>
                        {organizations.map((org) => (
                          <DropdownMenuItem
                            key={org.id}
                            onClick={() => onSwitchOrganization?.(org.id)}
                          >
                            <span>{org.name}</span>
                            {org.id === currentOrganizationId && (
                              <Badge
                                className="ml-auto"
                                variant="secondary"
                              >
                                Current
                              </Badge>
                            )}
                          </DropdownMenuItem>
                        ))}
                      </DropdownMenuSubContent>
                    </DropdownMenuPortal>
                  </DropdownMenuSub>
                  <DropdownMenuSeparator />
                </>
              )}
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={onNavigateToProfile}>
                <SparklesIcon className="mr-2 h-4 w-4" />
                Profile
                <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={onNavigateToBilling}>
                <CreditCardIcon className="mr-2 h-4 w-4" />
                Billing
                <DropdownMenuShortcut>⌘B</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={onLogout}>
              <LogOutIcon className="mr-2 h-4 w-4" />
              Log out
              <DropdownMenuShortcut>⌘Q</DropdownMenuShortcut>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
