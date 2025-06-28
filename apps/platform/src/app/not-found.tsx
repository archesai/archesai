'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'

import { Button } from '@archesai/ui/components/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

export default function NotFound() {
  const router = useRouter()

  return (
    <div className='flex min-h-screen items-center justify-center p-4'>
      <Card className='w-full max-w-md text-center'>
        <CardHeader>
          <div className='text-6xl font-bold text-muted-foreground'>404</div>
          <CardTitle className='text-2xl'>Page Not Found</CardTitle>
          <CardDescription>
            The page you're looking for doesn't exist or has been moved.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className='flex flex-col gap-2 sm:flex-row'>
            <Button
              className='flex-1'
              onClick={() => {
                router.back()
              }}
              variant='outline'
            >
              Go Back
            </Button>
            <Button
              asChild
              className='flex-1'
            >
              <Link href='/'>Go Home</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
