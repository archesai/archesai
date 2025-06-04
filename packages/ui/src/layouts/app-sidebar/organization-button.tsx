'use client'

import { ChevronsUpDown, Plus } from 'lucide-react'
import { toast } from 'sonner'

import type {
  MemberEntity,
  OrganizationEntity,
  UserEntity
} from '@archesai/domain'

import { findManyMembers, getOneUser, updateUser } from '@archesai/client'

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

export async function OrganizationButton({
  organization
}: OrganizationButtonProps) {
  const { data: user, status } = await getOneUser('organziation-button')
  if (status !== 200) {
    return null
  }
  const { data: memberships } = await findManyMembers({
    filter: {
      orgname: {
        equals: user.data.attributes.orgname
      }
    }
  })

  const { isMobile } = useSidebar()

  const handleSwitchOrganization = async (orgname: string) => {
    const { data, status } = await updateUser('organization-bar', {
      orgname
    })

    if (status === 200) {
      toast('Default organization updated', {
        description: 'Your default organization has been updated.'
      })
    } else {
      toast('Error updating default organization', {
        description: data.errors.map((error) => error.title).join(' ')
      })
    }
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
            {memberships.data.map((membership) => (
              <DropdownMenuItem
                className='flex justify-between gap-2'
                key={membership.id}
                onClick={async () =>
                  handleSwitchOrganization(membership.attributes.orgname)
                }
              >
                {membership.attributes.orgname}
                {user.data.attributes.orgname ===
                  membership.attributes.orgname && <Badge>Current</Badge>}
              </DropdownMenuItem>
            ))}
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
