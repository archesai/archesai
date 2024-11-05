"use client";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import { useThreadsControllerFindAll } from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import {
  Folder,
  Forward,
  ListMinus,
  MoreHorizontal,
  Trash2,
} from "lucide-react";
import Link from "next/link";

export function RecentThreads() {
  const { isMobile } = useSidebar();
  const { defaultOrgname } = useAuth();
  const { data: threads } = useThreadsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
  });

  return (
    <SidebarGroup className="group-data-[collapsible=icon]:hidden">
      <SidebarGroupLabel>Recent Threads</SidebarGroupLabel>
      <SidebarMenu>
        {threads?.results?.map((thread) => (
          <SidebarMenuItem key={thread.name}>
            <SidebarMenuButton asChild>
              <Link
                href={`/chatbots/single?chatbotId=${thread.chatbotId}&threadId=${thread.id}`}
              >
                <ListMinus />
                <span>{thread.name}</span>
              </Link>
            </SidebarMenuButton>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuAction showOnHover>
                  <MoreHorizontal />
                  <span className="sr-only">More</span>
                </SidebarMenuAction>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align={isMobile ? "end" : "start"}
                className="w-48 rounded-lg"
                side={isMobile ? "bottom" : "right"}
              >
                <DropdownMenuItem>
                  <Folder className="text-muted-foreground" />
                  <span>View Project</span>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Forward className="text-muted-foreground" />
                  <span>Share Project</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Trash2 className="text-muted-foreground" />
                  <span>Delete Project</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        ))}
        <SidebarMenuItem>
          <SidebarMenuButton className="text-sidebar-foreground/70">
            <MoreHorizontal className="text-sidebar-foreground/70" />
            <Link href="/chatbots/threads">More</Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroup>
  );
}
