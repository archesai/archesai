'use client'

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '@/components/ui/sidebar'
import { useLabelsControllerFindAll } from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import {
  Folder,
  Forward,
  ListMinus,
  MoreHorizontal,
  Trash2
} from 'lucide-react'
import Link from 'next/link'

export function RecentLabels() {
  const { isMobile } = useSidebar()
  const { defaultOrgname } = useAuth()
  const { data: labels } = useLabelsControllerFindAll({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  return (
    <SidebarGroup className='group-data-[collapsible=icon]:hidden'>
      <SidebarGroupLabel>Recent Labels</SidebarGroupLabel>
      <SidebarMenu>
        {labels?.results?.map((label) => (
          <SidebarMenuItem key={label.id}>
            <SidebarMenuButton asChild>
              <Link href={`/chatbots/chat?labelId=${label.id}`}>
                <ListMinus />
                <span>{label.name}</span>
              </Link>
            </SidebarMenuButton>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuAction showOnHover>
                  <MoreHorizontal />
                  <span className='sr-only'>More</span>
                </SidebarMenuAction>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align={isMobile ? 'end' : 'start'}
                className='w-48 rounded-lg'
                side={isMobile ? 'bottom' : 'right'}
              >
                <DropdownMenuItem>
                  <Folder className='text-muted-foreground' />
                  <span>View Project</span>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Forward className='text-muted-foreground' />
                  <span>Share Project</span>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Trash2 className='text-muted-foreground' />
                  <span>Delete Project</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        ))}
        <SidebarMenuItem>
          <SidebarMenuButton className='text-sidebar-foreground/70'>
            <MoreHorizontal className='text-sidebar-foreground/70' />
            <Link href='/chatbots/labels'>More</Link>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroup>
  )
}
