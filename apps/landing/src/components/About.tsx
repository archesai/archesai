import { ArchesLogo } from '@archesai/ui/components/custom/arches-logo'

import { Statistics } from './Statistics'

export const About = () => {
  return (
    <section
      className='container py-24 sm:py-32'
      id='about'
    >
      <div className='rounded-lg border bg-muted/50 py-12'>
        <div className='flex flex-col-reverse items-center gap-8 px-6 md:flex-row md:gap-12'>
          <ArchesLogo size='sm' />
          <div className='bg-green-0 flex flex-col justify-between'>
            <div className='pb-6'>
              <h2 className='text-3xl font-bold md:text-4xl'>
                <span className='bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent'>
                  About{' '}
                </span>
                Arches AI
              </h2>
              <p className='mt-4 text-xl text-muted-foreground'>
                At Arches AI, we are revolutionizing the way businesses
                integrate artificial intelligence into their workflows. Our
                platform empowers organizations with advanced AI capabilities,
                enabling them to automate processes, gain insights, and deliver
                better outcomes at scale. From AI-driven chatbots to robust data
                analysis tools, Arches AI provides tailored solutions to meet
                the needs of modern enterprises.
              </p>
            </div>

            <Statistics />
          </div>
        </div>
      </div>
    </section>
  )
}
