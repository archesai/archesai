import { ChevronsUpDown, Plus } from 'lucide-react'

import type { UserEntity } from '@archesai/schemas'

import {
  useFindManyMembersSuspense,
  useGetSessionSuspense,
  useUpdateUser
} from '@archesai/client'

import { ArchesLogo } from '#components/custom/arches-logo'
import { Badge } from '#components/shadcn/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '#components/shadcn/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '#components/shadcn/sidebar'

export function OrganizationButton() {
  const { data: sessionData } = useGetSessionSuspense()
  const { data: memberships } = useFindManyMembersSuspense({
    filter: {
      field: 'organizationId',
      operator: 'eq',
      type: 'condition',
      value: sessionData.session.activeOrganizationId ?? ''
    }
  })

  const { mutateAsync: updateUser } = useUpdateUser()

  const { isMobile } = useSidebar()

  const handleSwitchOrganization = async (_organizationId: string) => {
    await updateUser({
      data: {} as UserEntity,
      id: 'organization'
    })
  }

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              size='lg'
            >
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg text-sidebar-primary-foreground'>
                <div className='-mt-0.5'>
                  <ArchesLogo
                    scale={0.12}
                    size='sm'
                  />
                </div>
              </div>
              <div className='grid flex-1 text-left text-sm leading-tight'>
                <span className='truncate font-medium'>
                  {sessionData.user.name}
                </span>
                <span className='truncate text-xs'>
                  {sessionData.session.activeOrganizationId}
                </span>
              </div>
              <ChevronsUpDown className='ml-auto' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align='start'
            className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              Organizations
            </DropdownMenuLabel>
            {memberships.data.map((membership) => (
              <DropdownMenuItem
                className='gap-2 p-2'
                key={membership.id}
                onClick={async () => {
                  await handleSwitchOrganization(membership.organizationId)
                }}
              >
                {membership.organizationId}
                {'Arches Platform' === membership.organizationId && (
                  <Badge>Current</Badge>
                )}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className='gap-2 p-2'>
              <div className='flex size-6 items-center justify-center rounded-md border bg-transparent'>
                <Plus className='size-4' />
              </div>
              <div className='font-medium text-muted-foreground'>
                New Organization
              </div>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
