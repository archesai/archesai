'use client'

import type { FormEvent } from 'react'

import { Button } from '@archesai/ui/components/shadcn/button'
import { Input } from '@archesai/ui/components/shadcn/input'

export const Newsletter = () => {
  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
  }

  return (
    <section id='newsletter'>
      <hr className='mx-auto w-11/12' />

      <div className='container py-24 sm:py-32'>
        <h3 className='text-center text-4xl font-bold md:text-5xl'>
          Join Our Daily{' '}
          <span className='from-primary/60 to-primary bg-gradient-to-b bg-clip-text text-transparent'>
            Newsletter
          </span>
        </h3>
        <p className='text-muted-foreground mt-4 mb-8 text-center text-xl'>
          Lorem ipsum dolor sit amet consectetur.
        </p>

        <form
          className='mx-auto flex w-full flex-col gap-4 md:w-6/12 md:flex-row md:gap-2 lg:w-4/12'
          onSubmit={handleSubmit}
        >
          <Input
            aria-label='email'
            className='bg-muted/50 dark:bg-muted/80'
            placeholder='leomirandadev@gmail.com'
          />
          <Button>Subscribe</Button>
        </form>
      </div>

      <hr className='mx-auto w-11/12' />
    </section>
  )
}
