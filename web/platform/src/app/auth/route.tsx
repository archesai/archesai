import { createFileRoute, Link, Outlet, redirect } from '@tanstack/react-router'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'

export const Route = createFileRoute('/auth')({
  beforeLoad: ({ context }) => {
    const REDIRECT_URL = '/'
    if (context.session?.user) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error
      throw redirect({
        to: REDIRECT_URL
      })
    }
    return {
      redirectUrl: REDIRECT_URL
    }
  },
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  return (
    <div className='flex h-svh flex-col items-center justify-center gap-4'>
      <Link to='/'>
        <ArchesLogo size='lg' />
      </Link>
      <Outlet />
    </div>
  )
}
