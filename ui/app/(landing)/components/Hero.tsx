import { Button } from '@/components/ui/button'
import { buttonVariants } from '@/components/ui/button'
import { GitHubLogoIcon } from '@radix-ui/react-icons'
import Link from 'next/link'

export const Hero = () => {
  return (
    <section className='container z-10 mt-10 grid h-svh max-w-4xl place-items-center gap-10 py-20 md:py-32'>
      <div className='z-10 space-y-6 text-center'>
        <main className='text-5xl font-bold md:text-5xl'>
          Elevate Your Business with AI-Driven Solutions
        </main>

        <p className='mx-auto text-xl text-muted-foreground md:w-10/12'>
          Create intelligent chatbots, generate vibrant AI visuals, and
          integrate seamlessly using our API or no-code widgets.
        </p>

        <div className='space-y-4 md:space-x-4 md:space-y-0'>
          <Button className='w-full md:w-1/3'>
            <Link href='/playground'>Get Started</Link>
          </Button>

          <a
            className={`w-full md:w-1/3 ${buttonVariants({
              variant: 'outline'
            })}`}
            href='https://github.com/leoMirandaa/shadcn-landing-page.git'
            rel='noreferrer noopener'
            target='_blank'
          >
            Github Repository
            <GitHubLogoIcon className='ml-2 h-5 w-5' />
          </a>
        </div>
      </div>

      {/* Hero cards sections */}
      <div className='z-10'>{/* <HeroCards /> */}</div>

      {/* Shadow effect */}
      <div className='shadow'></div>
    </section>
  )
}
