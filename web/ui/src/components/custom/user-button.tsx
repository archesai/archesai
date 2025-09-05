import { useQueryClient } from '@tanstack/react-query'
import { useNavigate, useRouter } from '@tanstack/react-router'
import { toast } from 'sonner'

import type { UserEntity } from '#types/entities'

import {
  BadgeCheckIcon,
  ChevronsUpDownIcon,
  CreditCardIcon,
  LogOutIcon,
  SparklesIcon
} from '#components/custom/icons'
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

interface UserButtonProps {
  memberships?: {
    data: {
      id: string
      organizationId: string
    }[]
  }
  onLogout?: () => Promise<void>
  onUpdateUser?: (
    data: {
      data: UserEntity
      id: string
    },
    options?: {
      onSuccess?: () => void
    }
  ) => Promise<void>
  sessionData?: {
    session: {
      activeOrganizationId?: null | string
    }
    user: {
      email: string
      id: string
      image?: null | string
      name: string
    }
  }
  side?: 'bottom' | 'left' | 'right' | 'top'
  size?: 'default' | 'lg' | 'sm' | null | undefined
}

export function UserButton({
  memberships,
  onLogout,
  onUpdateUser,
  sessionData,
  size = 'lg'
}: UserButtonProps) {
  const defaultOrgname = 'Arches Platform'
  const { isMobile } = useSidebar()
  const navigation = useNavigate()
  const router = useRouter()
  const queryClient = useQueryClient()

  // If no session data is provided, render a placeholder
  if (!sessionData) {
    return (
      <SidebarMenu>
        <SidebarMenuItem>
          <Button
            className='h-8 w-8'
            size='sm'
            variant={'ghost'}
          >
            <Avatar>
              <AvatarFallback>U</AvatarFallback>
            </Avatar>
          </Button>
        </SidebarMenuItem>
      </SidebarMenu>
    )
  }

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
                <ChevronsUpDownIcon className='ml-auto size-4' />
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
                <BadgeCheckIcon />
                Profile
                <DropdownMenuShortcut>⌘P</DropdownMenuShortcut>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuSub>
              <DropdownMenuSubTrigger>Organizations</DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  {memberships?.data.map((membership) => (
                    <DropdownMenuItem
                      className='flex justify-between gap-2'
                      key={membership.id}
                      onClick={async () => {
                        if (onUpdateUser) {
                          await onUpdateUser(
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
                        }
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
                <SparklesIcon />
                Upgrade to Pro
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem>
                <CreditCardIcon />
                Billing
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={async () => {
                if (onLogout) {
                  await onLogout()
                }
                await navigation({
                  to: '/auth/login'
                })
                await router.invalidate()
                await queryClient.invalidateQueries()
                toast('Logged out successfully', {
                  description: 'You have been logged out.'
                })
              }}
            >
              <LogOutIcon />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
