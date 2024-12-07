import { buttonVariants } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { Facebook, Instagram, Linkedin } from 'lucide-react'
import Image from 'next/image'

interface SociaNetworkslProps {
  name: string
  url: string
}

interface TeamProps {
  imageUrl: string
  name: string
  position: string
  socialNetworks: SociaNetworkslProps[]
}

const teamList: TeamProps[] = [
  {
    imageUrl: 'https://i.pravatar.cc/150?img=35',
    name: 'Emma Smith',
    position: 'Product Manager',
    socialNetworks: [
      { name: 'Linkedin', url: 'http://linkedin.com' },
      {
        name: 'Facebook',
        url: 'https://www.facebook.com/'
      },
      {
        name: 'Instagram',
        url: 'https://www.instagram.com/'
      }
    ]
  },
  {
    imageUrl: 'https://i.pravatar.cc/150?img=60',
    name: 'John Doe',
    position: 'Tech Lead',
    socialNetworks: [
      { name: 'Linkedin', url: 'http://linkedin.com' },
      {
        name: 'Facebook',
        url: 'https://www.facebook.com/'
      },
      {
        name: 'Instagram',
        url: 'https://www.instagram.com/'
      }
    ]
  },
  {
    imageUrl: 'https://i.pravatar.cc/150?img=36',
    name: 'Ashley Ross',
    position: 'Frontend Developer',
    socialNetworks: [
      { name: 'Linkedin', url: 'http://linkedin.com' },

      {
        name: 'Instagram',
        url: 'https://www.instagram.com/'
      }
    ]
  },
  {
    imageUrl: 'https://i.pravatar.cc/150?img=17',
    name: 'Bruce Rogers',
    position: 'Backend Developer',
    socialNetworks: [
      { name: 'Linkedin', url: 'http://linkedin.com' },
      {
        name: 'Facebook',
        url: 'https://www.facebook.com/'
      }
    ]
  }
]

export const Team = () => {
  const socialIcon = (iconName: string) => {
    switch (iconName) {
      case 'Facebook':
        return <Facebook size='20' />

      case 'Instagram':
        return <Instagram size='20' />

      case 'Linkedin':
        return <Linkedin size='20' />
    }
  }

  return (
    <section
      className='container py-24 sm:py-32'
      id='team'
    >
      <h2 className='text-3xl font-bold md:text-4xl'>
        <span className='bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent'>
          Our Dedicated{' '}
        </span>
        Crew
      </h2>

      <p className='mb-10 mt-4 text-xl text-muted-foreground'>
        Lorem ipsum dolor sit amet consectetur, adipisicing elit. Veritatis
        dolor pariatur sit!
      </p>

      <div className='grid gap-8 gap-y-10 md:grid-cols-2 lg:grid-cols-4'>
        {teamList.map(
          ({ imageUrl, name, position, socialNetworks }: TeamProps) => (
            <Card
              className='relative mt-8 flex flex-col items-center justify-center bg-muted/50'
              key={name}
            >
              <CardHeader className='mt-8 flex items-center justify-center pb-2'>
                <Image
                  alt={`${name} ${position}`}
                  className='absolute -top-12 aspect-square h-24 w-24 rounded-full object-cover'
                  height={96}
                  src={imageUrl}
                  width={96}
                />
                <CardTitle className='text-center'>{name}</CardTitle>
                <CardDescription className='text-primary'>
                  {position}
                </CardDescription>
              </CardHeader>

              <CardContent className='pb-2 text-center'>
                <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
              </CardContent>

              <CardFooter>
                {socialNetworks.map(({ name, url }: SociaNetworkslProps) => (
                  <div key={name}>
                    <a
                      className={buttonVariants({
                        size: 'sm',
                        variant: 'ghost'
                      })}
                      href={url}
                      rel='noreferrer noopener'
                      target='_blank'
                    >
                      <span className='sr-only'>{name} icon</span>
                      {socialIcon(name)}
                    </a>
                  </div>
                ))}
              </CardFooter>
            </Card>
          )
        )}
      </div>
    </section>
  )
}
