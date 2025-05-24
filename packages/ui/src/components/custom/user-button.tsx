'use client'

import { Link } from '@radix-ui/react-navigation-menu'
import {
  BadgeCheck,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles
} from 'lucide-react'
import { toast } from 'sonner'

import { updateUser, useFindManyMembers, useGetOneUser } from '@archesai/client'

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
import { useAuth } from '#hooks/use-auth'
import { cn } from '#lib/utils'

export function UserButton({
  side = 'right',
  size = 'lg'
}: {
  side?: 'bottom' | 'left' | 'right' | 'top'
  size?: 'default' | 'lg' | 'sm' | null | undefined
}) {
  const { isMobile } = useSidebar()
  const { defaultOrgname, logout } = useAuth()
  const { data, status } = useGetOneUser('me')
  if (status === 'pending') return null
  if (data?.status === 404) {
    return null
  }
  const user = data?.data.data
  if (!user) {
    return null
  }

  const { data: memberships } = useFindManyMembers({
    filter: {
      userId: {
        equals: user.id
      }
    }
  })
  if (!memberships) {
    return null
  }

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
                  alt={user.attributes.name}
                  src={user.attributes.image ?? ''}
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
                    alt={user.attributes.name}
                    src={user.attributes.image ?? ''}
                  />
                  <AvatarFallback className='rounded-lg'>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>
                    {user.attributes.name}
                  </span>
                  <span className='truncate text-xs'>
                    {user.attributes.email}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Link href='/profile/general'>
                  <BadgeCheck />
                  Profile
                  <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
                </Link>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Organizations</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  {memberships.data.data.map((membership) => (
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
                  ))}
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem>
              <Link href='/organization/general'>Settings</Link>
              <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Link href='/organization/billing'>
                  <Sparkles />
                  Upgrade to Pro
                </Link>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Link href='/organization/billing'>
                  <CreditCard />
                  Billing
                </Link>
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
