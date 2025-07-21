import { createFileRoute, Outlet } from '@tanstack/react-router'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'
import { useLinkComponent } from '@archesai/ui/hooks/use-link'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  const Link = useLinkComponent()
  return (
    <div className='flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10'>
      <div className='flex w-full max-w-sm flex-col gap-6'>
        <Link
          className='flex items-center gap-2 self-center font-medium'
          href='/'
        >
          <ArchesLogo size='lg' />
        </Link>
        <Outlet />
      </div>
    </div>
  )
}
