import { createFileRoute, Outlet } from '@tanstack/react-router'

import { CpuIcon } from '@archesai/ui/components/custom/icons'
import { useLinkComponent } from '@archesai/ui/hooks/use-link'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  const Link = useLinkComponent()
  return (
    <div className='flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10'>
      <div className='flex w-full max-w-sm flex-col gap-6'>
        <Link
          className='flex items-center gap-2 self-center font-medium'
          href='/'
        >
          <div className='flex size-6 items-center justify-center rounded-md bg-primary text-primary-foreground'>
            <CpuIcon className='size-4' />
          </div>
          Acme Inc.
        </Link>
        <Outlet />
      </div>
    </div>
  )
}
