import type { JSX } from "react";

import { toast } from "sonner";

import { ArchesLogo } from "#components/custom/arches-logo";
import { ChevronsUpDownIcon, PlusSquareIcon } from "#components/custom/icons";
import { Badge } from "#components/shadcn/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#components/shadcn/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "#components/shadcn/sidebar";

interface OrganizationButtonProps {
  memberships?: {
    data: {
      id: string;
      organizationID: string;
    }[];
  };
  onUpdateSession?: (
    data: {
      data: {
        activeOrganizationID: string;
      };
      id: string;
    },
    options?: {
      onSuccess?: () => void;
    },
  ) => Promise<void>;
  session?: {
    activeOrganizationID?: null | string;
    id: string;
  };
  user?: {
    email: string;
    id: string;
    name: string;
  };
}

export function OrganizationButton({
  memberships,
  onUpdateSession,
  session,
}: OrganizationButtonProps): JSX.Element {
  const { isMobile } = useSidebar();

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              <ArchesLogo
                scale={0.8}
                size="sm"
              />
              <ChevronsUpDownIcon className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-muted-foreground text-xs">
              Organizations
            </DropdownMenuLabel>
            {memberships?.data.map((membership) => (
              <DropdownMenuItem
                className="gap-2 p-2"
                key={membership.id}
                onClick={async () => {
                  if (onUpdateSession && session) {
                    await onUpdateSession(
                      {
                        data: {
                          activeOrganizationID: membership.organizationID,
                        },
                        id: session.id,
                      },
                      {
                        onSuccess: () => {
                          toast.success(
                            `Switched to organization: ${membership.organizationID}`,
                          );
                        },
                      },
                    );
                  }
                }}
              >
                {membership.organizationID}
                {session?.activeOrganizationID ===
                  membership.organizationID && <Badge>Current</Badge>}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="gap-2 p-2">
              <div className="flex size-6 items-center justify-center rounded-md border bg-transparent">
                <PlusSquareIcon className="size-4" />
              </div>
              <div className="font-medium text-muted-foreground">
                New Organization
              </div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
