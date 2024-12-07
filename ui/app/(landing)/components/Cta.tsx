import { Button } from '@/components/ui/button'

export const Cta = () => {
  return (
    <section className='my-24 bg-muted/50 py-16 sm:my-32' id='cta'>
      <div className='container place-items-center lg:grid lg:grid-cols-2'>
        <div className='lg:col-start-1'>
          <h2 className='text-3xl font-bold md:text-4xl'>
            All Your
            <span className='bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent'>
              {' '}
              Ideas & Concepts{' '}
            </span>
            In One Interface
          </h2>
          <p className='mb-8 mt-4 text-xl text-muted-foreground lg:mb-0'>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Eaque, beatae. Ipsa tempore ipsum iste quibusdam
            illum ducimus eos. Quasi, sed!
          </p>
        </div>

        <div className='space-y-4 lg:col-start-2'>
          <Button className='w-full md:mr-4 md:w-auto'>Request a Demo</Button>
          <Button className='w-full md:w-auto' variant='outline'>
            View all features
          </Button>
        </div>
      </div>
    </section>
  )
}
