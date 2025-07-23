import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'
import { useLinkComponent } from '@archesai/ui/hooks/use-link'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout,
  beforeLoad: async ({ context }) => {
    const REDIRECT_URL = '/'
    if (context.session?.user) {
      throw redirect({
        to: REDIRECT_URL
      })
    }
    return {
      redirectUrl: REDIRECT_URL
    }
  }
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
