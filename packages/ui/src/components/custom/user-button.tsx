import {
  BadgeCheck,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles
} from 'lucide-react'
import { toast } from 'sonner'

import {
  useFindManyMembers,
  useGetSession,
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
import { Skeleton } from '#components/shadcn/skeleton'

export function UserButton({
  size = 'lg'
}: {
  side?: 'bottom' | 'left' | 'right' | 'top'
  size?: 'default' | 'lg' | 'sm' | null | undefined
}) {
  const defaultOrgname = 'Arches Platform'
  const { data: user } = useGetSession({
    fetch: {
      credentials: 'include'
    }
  })
  const { mutateAsync: logout } = useLogout({
    fetch: {
      credentials: 'include',
      redirect: 'error'
    }
  })
  const { isMobile } = useSidebar()

  const { data: memberships } = useFindManyMembers(
    {
      filter: {
        userId: {
          equals: 'Arches Platform'
        }
      }
    },
    {
      query: {
        enabled: false,
        initialData: {
          data: []
        }
      }
    }
  )

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
                    alt={user?.name}
                    src={user?.image}
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
                    alt={user?.name}
                    src={user?.image}
                  />
                  <AvatarFallback>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>user.name</span>
                  <span className='truncate text-xs'>user.email</span>
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
                    alt={user?.name}
                    src={user?.image}
                  />
                  <AvatarFallback className='rounded-lg'>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-medium'>
                    {user ? user.name : <Skeleton />}
                  </span>
                  <span className='truncate text-xs'>
                    {user ? user.email : <Skeleton />}
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
                            data: {
                              orgname: membership.attributes.orgname
                            },
                            id: ''
                          },
                          {
                            onSuccess: () => {
                              toast('Organization changed', {
                                description: `You have
                              switched to ${membership.attributes.orgname}`
                              })
                            }
                          }
                        )
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
