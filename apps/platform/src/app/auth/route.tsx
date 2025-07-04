import { createFileRoute, Outlet } from '@tanstack/react-router'

import { CpuIcon } from '@archesai/ui/components/custom/icons'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  return (
    <div className='flex min-h-svh flex-col items-center justify-center gap-6 bg-muted p-6 md:p-10'>
      <div className='flex w-full max-w-sm flex-col gap-6'>
        <a
          className='flex items-center gap-2 self-center font-medium'
          href='#'
        >
          <div className='flex size-6 items-center justify-center rounded-md bg-primary text-primary-foreground'>
            <CpuIcon className='size-4' />
          </div>
          Acme Inc.
        </a>
        <Outlet />
      </div>
    </div>
  )
}
