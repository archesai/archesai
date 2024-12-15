'use client'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
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
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar
} from '@/components/ui/sidebar'
import { Skeleton } from '@/components/ui/skeleton'
import {
  useUsersControllerFindOne,
  useUsersControllerUpdate
} from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { useToast } from '@/hooks/use-toast'
import { cn } from '@/lib/utils'
import {
  BadgeCheck,
  ChevronsUpDown,
  CreditCard,
  LogOut,
  Sparkles
} from 'lucide-react'
import { useRouter } from 'next/navigation'

export function UserButton({
  size = 'lg',
  side = 'right'
}: {
  size?: 'lg' | 'sm' | 'default' | undefined | null
  side?: 'left' | 'right' | 'top' | 'bottom'
}) {
  const router = useRouter()
  const { toast } = useToast()
  const { isMobile } = useSidebar()
  const { data: user, isFetched } = useUsersControllerFindOne({})
  const { logout, defaultOrgname } = useAuth()
  const { mutateAsync: updateDefaultOrg } = useUsersControllerUpdate({
    onError: (error) => {
      toast({
        description: error?.message,
        title: 'Error updating default organization',
        variant: 'destructive'
      })
    },
    onSuccess: () => {
      toast({
        description: 'Your default organization has been updated.',
        title: 'Default organization updated'
      })
    }
  })
  const memberships = user?.memberships

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
                  alt={user?.displayName}
                  src={user?.photoUrl}
                />
                <AvatarFallback className='rounded-xl'>
                  <Skeleton className='h-9 w-9 bg-sidebar-accent' />
                </AvatarFallback>
              </Avatar>
              {size !== 'sm' && (
                <>
                  <div className='grid flex-1 text-left text-sm leading-tight'>
                    <span className='truncate font-semibold'>
                      {isFetched ? (
                        user?.displayName
                      ) : (
                        <Skeleton className='m-1 h-4 bg-sidebar-accent' />
                      )}
                    </span>
                    <span className='truncate text-xs'>
                      {isFetched ? (
                        user?.email
                      ) : (
                        <Skeleton className='m-1 h-3 bg-sidebar-accent' />
                      )}
                    </span>
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
                    alt={user?.displayName}
                    src={user?.photoUrl}
                  />
                  <AvatarFallback className='rounded-lg'>CN</AvatarFallback>
                </Avatar>
                <div className='grid flex-1 text-left text-sm leading-tight'>
                  <span className='truncate font-semibold'>
                    {user?.displayName}
                  </span>
                  <span className='truncate text-xs'>{user?.email}</span>
                </div>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem onClick={() => router.push('/profile/general')}>
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
                  {memberships?.map((membership) => (
                    <DropdownMenuItem
                      className='flex justify-between gap-2'
                      key={membership.id}
                      onClick={() => {
                        updateDefaultOrg({
                          body: {
                            defaultOrgname: membership.orgname
                          }
                        })
                      }}
                    >
                      {membership.orgname}
                      {defaultOrgname === membership.orgname && (
                        <Badge>Current</Badge>
                      )}
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuItem
              onClick={() => router.push('/organization/general')}
            >
              Settings
              <DropdownMenuShortcut>⌘S</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem
                onClick={() => router.push('/organization/billing')}
              >
                <Sparkles />
                Upgrade to Pro
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuGroup>
              <DropdownMenuItem
                onClick={() => router.push('/organization/billing')}
              >
                <CreditCard />
                Billing
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={async () => await logout()}>
              <LogOut />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
