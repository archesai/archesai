import type { ErrorComponentProps } from '@tanstack/react-router'

import { useQueryClient } from '@tanstack/react-query'
import { Link, rootRouteId, useMatch, useRouter } from '@tanstack/react-router'

import { Button } from '@archesai/ui/components/shadcn/button'

export function DefaultCatchBoundary({ error }: ErrorComponentProps) {
  const router = useRouter()
  const queryClient = useQueryClient()
  const isRoot = useMatch({
    select: (state) => state.id === rootRouteId,
    strict: false
  })

  return (
    <div className='z-0 mb-16 flex h-full flex-col items-center justify-center gap-6 bg-primary/0 dark:bg-primary/0'>
      <h1 className='text-2xl font-bold'>Something went wrong</h1>
      <p className='text-destructive'>
        {error.message || 'An unexpected error occurred.'}
      </p>
      <div className='flex flex-wrap items-center gap-2'>
        <Button
          onClick={async () => {
            await router.invalidate()
            await queryClient.invalidateQueries()
          }}
          size='sm'
          variant={'ghost'}
        >
          Try Again
        </Button>
        {isRoot ?
          <Link to='/'>Home</Link>
        : <Button
            asChild
            size='sm'
            variant={'ghost'}
          >
            <Link
              onClick={(e) => {
                e.preventDefault()
                window.history.back()
              }}
              to='/'
            >
              Go Back
            </Link>
          </Button>
        }
      </div>
    </div>
  )
}
