import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'

interface TestimonialProps {
  comment: string
  image: string
  name: string
  userName: string
}

const testimonials: TestimonialProps[] = [
  {
    comment: 'This landing page is awesome!',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe'
  },
  {
    comment:
      'Lorem ipsum dolor sit amet,empor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud.',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe1'
  },

  {
    comment:
      'Lorem ipsum dolor sit amet,exercitation. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe2'
  },
  {
    comment:
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam.',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe3'
  },
  {
    comment:
      'Lorem ipsum dolor sit amet, tempor incididunt  aliqua. Ut enim ad minim veniam, quis nostrud.',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe4'
  },
  {
    comment:
      'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
    image: 'https://github.com/shadcn.png',
    name: 'John Doe React',
    userName: '@john_Doe5'
  }
]

export const Testimonials = () => {
  return (
    <section
      className='container py-24 sm:py-32'
      id='testimonials'
    >
      <h2 className='text-3xl font-bold md:text-4xl'>
        Discover Why
        <span className='from-primary/60 to-primary bg-gradient-to-b bg-clip-text text-transparent'>
          {' '}
          People Love{' '}
        </span>
        This Landing Page
      </h2>

      <p className='text-muted-foreground pb-8 pt-4 text-xl'>
        Lorem ipsum dolor sit amet, consectetur adipisicing elit. Non unde error
        facere hic reiciendis illo
      </p>

      <div className='mx-auto grid columns-2 space-y-4 sm:block md:grid-cols-2 lg:columns-3 lg:grid-cols-4 lg:gap-6 lg:space-y-6'>
        {testimonials.map(
          ({ comment, image, name, userName }: TestimonialProps) => (
            <Card
              className='max-w-md overflow-hidden md:break-inside-avoid'
              key={userName}
            >
              <CardHeader className='flex flex-row items-center gap-4 pb-2'>
                <Avatar>
                  <AvatarImage
                    alt=''
                    src={image}
                  />
                  <AvatarFallback>OM</AvatarFallback>
                </Avatar>

                <div className='flex flex-col'>
                  <CardTitle className='text-lg'>{name}</CardTitle>
                  <CardDescription>{userName}</CardDescription>
                </div>
              </CardHeader>

              <CardContent>{comment}</CardContent>
            </Card>
          )
        )}
      </div>
    </section>
  )
}
