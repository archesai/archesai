'use client'

import {
  BadgeCheck,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles
} from 'lucide-react'
import { toast } from 'sonner'

import type { UserEntity } from '@archesai/domain'

import { logout, updateUser, useFindManyMembers } from '@archesai/client'

import { Avatar, AvatarFallback, AvatarImage } from '#components/shadcn/avatar'
import { Badge } from '#components/shadcn/badge'
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
  DropdownMenuTrigger
} from '#components/shadcn/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '#components/shadcn/sidebar'
import { Skeleton } from '#components/shadcn/skeleton'
import { cn } from '#lib/utils'

export function UserButton({
  side = 'right',
  size = 'lg'
}: {
  side?: 'bottom' | 'left' | 'right' | 'top'
  size?: 'default' | 'lg' | 'sm' | null | undefined
}) {
  const { isMobile } = useSidebar()
  const defaultOrgname = 'Arches Platform'

  const user = {} as UserEntity
  const userId = ''

  const { data: memberships } = useFindManyMembers(
    {
      filter: {
        userId: {
          equals: userId
        }
      }
    },
    {
      fetch: {
        credentials: 'include'
      },
      query: {
        enabled: !!userId
      }
    }
  )

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              className={cn(
                'data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground',
                size === 'sm' ? 'p-0' : 'p-2'
              )}
              size={size}
            >
              <Avatar className='h-8 w-8 rounded-lg'>
                <AvatarImage
                  alt={user.name}
                  src={user.image ?? ''}
                />
                <AvatarFallback className='rounded-xl'>
                  <Skeleton className='h-9 w-9 bg-sidebar-accent' />
                </AvatarFallback>
              </Avatar>
              {size !== 'sm' && (
                <>
                  <div className='grid flex-1 text-left text-sm leading-tight'>
                    <span className='truncate font-semibold'>user.name</span>
                    <span className='truncate text-xs'>user.email</span>
                  </div>
                  <ChevronsUpDown className='ml-auto size-4' />
                </>
              )}
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align='end'
            className='w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg'
            side={isMobile ? 'bottom' : side}
            sideOffset={4}
          >
            <DropdownMenuLabel className='p-0 font-normal'>
              <div className='flex items-center gap-2 px-1 py-1.5 text-left text-sm'>
                <Avatar className='h-8 w-8 rounded-lg'>
                  <AvatarImage
                    alt={user.name}
                    src={user.image ?? ''}
                  />
                  <AvatarFallback className='rounded-lg'>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>{user.name}</span>
                  <span className='truncate text-xs'>{user.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <a
                  className='flex w-full'
                  href='/profile/general'
                >
                  <BadgeCheck />
                  Profile
                  <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
                </a>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Organizations</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  {memberships ?
                    // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
                    (memberships.data.data ?? []).map((membership) => (
                      <DropdownMenuItem
                        className='flex justify-between gap-2'
                        key={membership.id}
                        onClick={async () => {
                          const { status } = await updateUser(user.id, {
                            orgname: membership.attributes.orgname
                          })

                          if (status === 200) {
                            toast('Organization changed', {
                              description: `You have
                              switched to ${membership.attributes.orgname}`
                            })
                          }
                        }}
                      >
                        {membership.attributes.orgname}
                        {defaultOrgname === membership.attributes.orgname && (
                          <Badge>Current</Badge>
                        )}
                      </DropdownMenuItem>
                    ))
                  : null}
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem>
              <a
                className='flex w-full'
                href='/organization/general'
              >
                Settings
              </a>
              <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <a
                  className='flex w-full'
                  href='/organization/billing'
                >
                  <Sparkles />
                  Upgrade to Pro
                </a>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <a
                  className='flex w-full'
                  href='/organization/billing'
                >
                  <CreditCard />
                  Billing
                </a>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={async () => {
                await logout()
              }}
            >
              <LogOut />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
