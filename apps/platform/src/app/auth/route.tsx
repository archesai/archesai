import { createFileRoute, Link, Outlet } from '@tanstack/react-router'

import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'
import { buttonVariants } from '@archesai/ui/components/shadcn/button'
import { cn } from '@archesai/ui/lib/utils'

export const Route = createFileRoute('/auth')({
  component: AuthenticationLayout
})

export default function AuthenticationLayout() {
  return (
    <div className='relative grid h-svh items-center justify-center lg:grid-cols-2'>
      <Link
        className={cn(
          buttonVariants({ variant: 'ghost' }),
          'absolute top-4 left-4 lg:right-4 lg:left-auto'
        )}
        to='/'
      >
        Back
      </Link>

      {/* Left side of the screen */}
      <div className='hidden h-full flex-col justify-between bg-primary p-10 text-white lg:flex'>
        <ArchesLogo />
        <blockquote className='flex flex-col gap-2'>
          <p className='text-lg'>
            &ldquo;This library has saved me countless hours of work and helped
            me deliver stunning designs to my clients faster than ever
            before.&rdquo;
          </p>
          <footer className='text-sm'>Sofia Davis</footer>
        </blockquote>
      </div>

      {/* Right side of the screen or main*/}
      <div className='mx-auto flex h-full max-w-xs flex-col justify-center gap-2'>
        <Outlet />
        <p className='text-center text-sm text-muted-foreground'>
          By clicking continue, you agree to our{' '}
          <Link
            className='underline underline-offset-4 hover:text-foreground'
            to='/legal/terms'
          >
            Terms of Service
          </Link>{' '}
          and{' '}
          <Link
            className='underline underline-offset-4 hover:text-foreground'
            to='/legal/privacy'
          >
            Privacy Policy
          </Link>
          .
        </p>
      </div>
    </div>
  )
}
