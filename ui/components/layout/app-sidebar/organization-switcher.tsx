"use client";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  useOrganizationsControllerFindOne,
  useUserControllerFindOne,
  useUserControllerUpdate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/use-auth";
import { ChevronsUpDown, Plus } from "lucide-react";
import * as React from "react";

import { ArchesLogo } from "../../arches-logo";
import { Badge } from "../../ui/badge";
import { useToast } from "../../ui/use-toast";

export function OrganizationSwitcher() {
  const { defaultOrgname } = useAuth();
  const { data: user } = useUserControllerFindOne({});
  const { data: organization } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname,
    },
  });
  const { toast } = useToast();

  const { isMobile } = useSidebar();
  const { mutateAsync: updateDefaultOrg } = useUserControllerUpdate({
    onError: (error) => {
      toast({
        description: error?.stack.message,
        title: "Error updating default organization",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "Your default organization has been updated.",
        title: "Default organization updated",
      });
    },
  });

  const memberships = user?.memberships;

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              size="lg"
            >
              {
                // add eslint ignore for the next lint
              }
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <div className="-mt-0.5">
                  <ArchesLogo scale={0.1} size="sm" />
                </div>
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">
                  {organization?.orgname}
                </span>
                <span className="truncate text-xs capitalize">
                  {organization?.plan?.toLocaleLowerCase()} Plan
                </span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="start"
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-xs text-muted-foreground">
              Teams
            </DropdownMenuLabel>
            {memberships?.map((membership) => (
              <DropdownMenuItem
                className="flex justify-between gap-2"
                key={membership.id}
                onClick={() => {
                  updateDefaultOrg({
                    body: {
                      defaultOrgname: membership.orgname,
                    },
                  });
                }}
              >
                {membership.orgname}
                {defaultOrgname === membership.orgname && (
                  <Badge>Current</Badge>
                )}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className="gap-2 p-2">
              <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                <Plus className="size-4" />
              </div>
              <div className="font-medium text-muted-foreground">Add team</div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
