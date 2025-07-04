import type { ErrorComponentProps } from '@tanstack/react-router'

import {
  ErrorComponent,
  Link,
  rootRouteId,
  useMatch,
  useRouter
} from '@tanstack/react-router'

import { Button } from '@archesai/ui/components/shadcn/button'

export function DefaultCatchBoundary({ error }: ErrorComponentProps) {
  const router = useRouter()
  const isRoot = useMatch({
    select: (state) => state.id === rootRouteId,
    strict: false
  })

  console.error(error)

  return (
    <div className='flex min-w-0 flex-1 flex-col items-center justify-center gap-6 p-4'>
      <ErrorComponent error={error} />
      <div className='flex flex-wrap items-center gap-2'>
        <Button
          onClick={async () => {
            await router.invalidate()
          }}
          size='sm'
          variant={'ghost'}
        >
          Try Again
        </Button>
        {isRoot ?
          <Link to='/chat'>Home</Link>
        : <Link
            onClick={(e) => {
              e.preventDefault()
              window.history.back()
            }}
            to='/chat'
          >
            Go Back
          </Link>
        }
      </div>
    </div>
  )
}
