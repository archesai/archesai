import { ChevronsUpDown, Plus } from 'lucide-react'
import { toast } from 'sonner'

import type {
  MemberEntity,
  OrganizationEntity,
  UserEntity
} from '@archesai/domain'

import {
  useFindManyMembers,
  useGetSession,
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
import { Skeleton } from '#components/shadcn/skeleton'

export interface OrganizationButtonProps {
  children?: never
  memberships?: MemberEntity[]
  organization?: OrganizationEntity
  user: UserEntity
}

export function OrganizationButton({ organization }: OrganizationButtonProps) {
  const { data: _session } = useGetSession({
    fetch: {
      credentials: 'include'
    }
  })

  const { data: memberships } = useFindManyMembers({
    filter: {
      orgname: {
        equals: 'Arches Platform'
      }
    }
  })

  const { mutateAsync: updateUser } = useUpdateUser()

  const { isMobile } = useSidebar()

  const handleSwitchOrganization = async (orgname: string) => {
    await updateUser(
      {
        data: {
          orgname
        },
        id: 'organization'
      },
      {
        onError: (errors) => {
          toast('Error updating default organization', {
            description: errors.errors.map((error) => error.title).join(' ')
          })
        },
        onSuccess: () => {
          toast('Default organization updated', {
            description: 'Your default organization has been updated.'
          })
        }
      }
    )
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
              <div className='flex aspect-square size-8 items-center justify-center rounded-lg border-sidebar-accent text-sidebar-foreground'>
                <div className='-mt-0.5'>
                  <ArchesLogo
                    scale={0.1}
                    size='sm'
                  />
                </div>
              </div>
              <div className='grid flex-1 text-left text-sm leading-tight'>
                <span className='truncate font-semibold'>
                  {!organization ?
                    <Skeleton className='m-1 h-2 bg-sidebar-accent' />
                  : organization.orgname}
                </span>
                <span className='truncate text-xs capitalize'>
                  {!organization ?
                    <Skeleton className='m-1 h-2 bg-sidebar-accent' />
                  : organization.plan + ' Plan'}
                </span>
              </div>
              <ChevronsUpDown className='ml-auto' />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align='start'
            className='w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='text-xs text-muted-foreground'>
              Organizations
            </DropdownMenuLabel>
            {memberships ?
              memberships.data.map((membership) => (
                <DropdownMenuItem
                  className='flex justify-between gap-2'
                  key={membership.id}
                  onClick={async () => {
                    await handleSwitchOrganization(
                      membership.attributes.orgname
                    )
                  }}
                >
                  {membership.attributes.orgname}
                  {'Arches Platform' === membership.attributes.orgname && (
                    <Badge>Current</Badge>
                  )}
                </DropdownMenuItem>
              ))
            : null}
            <DropdownMenuSeparator />
            <DropdownMenuItem className='gap-2 p-2'>
              <div className='flex size-6 items-center justify-center rounded-md border bg-background'>
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
