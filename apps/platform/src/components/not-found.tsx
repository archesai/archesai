import { Link } from '@tanstack/react-router'

import { Button } from '@archesai/ui/components/shadcn/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

export default function NotFound({ children }: { children?: React.ReactNode }) {
  return (
    <div className='flex min-h-screen items-center justify-center p-4'>
      <Card className='w-full max-w-md text-center'>
        <CardHeader>
          <div className='text-6xl font-bold text-muted-foreground'>404</div>
          <CardTitle className='text-2xl'>Page Not Found</CardTitle>
          <CardDescription>
            The page you&apos;re looking for doesn&apos;t exist or has been
            moved.
            {children}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className='flex flex-col gap-2 sm:flex-row'>
            <Button
              asChild
              className='flex-1'
            >
              <Link to='/chat'>Go Home</Link>
            </Button>
            <Button
              asChild
              className='flex-1'
            >
              <Link to='/chat'>Go Home</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
