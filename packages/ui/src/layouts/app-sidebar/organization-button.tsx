import { toast } from 'sonner'

import {
  useFindManyMembersSuspense,
  useGetSessionSuspense,
  useUpdateSession
} from '@archesai/client'

import { ArchesLogo } from '#components/custom/arches-logo'
import { ChevronsUpDownIcon, PlusSquareIcon } from '#components/custom/icons'
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
  const {
    data: { session, user }
  } = useGetSessionSuspense()
  const { data: memberships } = useFindManyMembersSuspense({
    filter: {
      field: 'userId',
      operator: 'eq',
      type: 'condition',
      value: user.id
    }
  })

  const { mutateAsync: updateSession } = useUpdateSession()

  const { isMobile } = useSidebar()

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className='data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
              size='lg'
            >
              <ArchesLogo
                scale={0.8}
                size='lg'
              />
              <ChevronsUpDownIcon className='ml-auto' />
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
                  await updateSession(
                    {
                      data: {
                        activeOrganizationId: membership.organizationId
                      },
                      id: session.id
                    },
                    {
                      onSuccess: () => {
                        toast.success(
                          `Switched to organization: ${membership.organizationId}`
                        )
                      }
                    }
                  )
                }}
              >
                {membership.organizationId}
                {session.activeOrganizationId === membership.organizationId && (
                  <Badge>Current</Badge>
                )}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <DropdownMenuItem className='gap-2 p-2'>
              <div className='flex size-6 items-center justify-center rounded-md border bg-transparent'>
                <PlusSquareIcon className='size-4' />
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
