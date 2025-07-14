import {
  BadgeCheck,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles
} from 'lucide-react'
import { toast } from 'sonner'

import type { UserEntity } from '@archesai/schemas'

import {
  useFindManyMembersSuspense,
  useGetSessionSuspense,
  useLogout,
  useUpdateUser
} from '@archesai/client'

import { Avatar, AvatarFallback, AvatarImage } from '#components/shadcn/avatar'
import { Badge } from '#components/shadcn/badge'
import { Button } from '#components/shadcn/button'
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

export function UserButton({
  size = 'lg'
}: {
  side?: 'bottom' | 'left' | 'right' | 'top'
  size?: 'default' | 'lg' | 'sm' | null | undefined
}) {
  const defaultOrgname = 'Arches Platform'
  const { data: sessionData } = useGetSessionSuspense()
  const { mutateAsync: logout } = useLogout()
  const { isMobile } = useSidebar()

  const { data: memberships } = useFindManyMembersSuspense({
    filter: {
      userId: {
        equals: 'Arches Platform'
      }
    }
  })

  const { mutateAsync: updateUser } = useUpdateUser()

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger
            asChild
            className='align-middle'
          >
            {size === 'sm' ?
              <Button
                className='h-8 w-8'
                size='sm'
                variant={'ghost'}
              >
                <Avatar>
                  <AvatarImage
                    alt={sessionData.user.email}
                    src={sessionData.user.image ?? undefined}
                  />
                  <AvatarFallback>CN</AvatarFallback>
                </Avatar>
              </Button>
            : <SidebarMenuButton
                className={
                  'data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground'
                }
                size={size}
              >
                <Avatar>
                  <AvatarImage
                    alt={sessionData.user.email}
                    src={sessionData.user.image ?? undefined}
                  />
                  <AvatarFallback>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>
                    {sessionData.user.name}
                  </span>
                  <span className='truncate text-xs'>
                    {sessionData.user.email}
                  </span>
                </div>
                <ChevronsUpDown className='ml-auto size-4' />
              </SidebarMenuButton>
            }
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align='end'
            className='w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg'
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className='p-0 font-normal'>
              <div className='flex items-center gap-2 px-1 py-1.5 text-left text-sm'>
                <Avatar className='h-8 w-8 rounded-lg'>
                  <AvatarImage
                    alt={sessionData.user.email}
                    src={sessionData.user.image ?? undefined}
                  />
                  <AvatarFallback className='rounded-lg'>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-medium'>
                    {sessionData.user.email}
                  </span>
                  <span className='truncate text-xs'>
                    {sessionData.user.email}
                  </span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <BadgeCheck />
                Profile
                <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Organizations</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  {memberships.data.map((membership) => (
                    <DropdownMenuItem
                      className='flex justify-between gap-2'
                      key={membership.id}
                      onClick={async () => {
                        await updateUser(
                          {
                            data: {} as UserEntity,
                            id: sessionData.user.id
                          },
                          {
                            onSuccess: () => {
                              toast('Organization changed', {
                                description: `You have
                              switched to ${membership.organizationId}`
                              })
                            }
                          }
                        )
                      }}
                    >
                      {membership.organizationId}
                      {defaultOrgname === membership.organizationId && (
                        <Badge>Current</Badge>
                      )}
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem>
              Settings
              <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <Sparkles />
                Upgrade to Pro
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <CreditCard />
                Billing
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={async () => {
                await logout()
                toast('Logged out successfully', {
                  description: 'You have been logged out.'
                })
                window.location.href = '/'
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
