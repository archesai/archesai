import { createFileRoute, Outlet } from '@tanstack/react-router'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'
import { useLinkComponent } from '@archesai/ui/hooks/use-link'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  const Link = useLinkComponent()
  return (
    <div className='flex h-svh flex-col items-center justify-center gap-4'>
      <Link href='/'>
        <ArchesLogo size='lg' />
      </Link>
      <Outlet />
    </div>
  )
}
