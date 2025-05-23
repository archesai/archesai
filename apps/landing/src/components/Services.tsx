import type { JSX } from 'react'

import Image from 'next/image'

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'

import { ChartIcon, MagnifierIcon, WalletIcon } from './Icons'

interface ServiceProps {
  description: string
  icon: JSX.Element
  title: string
}

const serviceList: ServiceProps[] = [
  {
    description:
      'Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.',
    icon: <ChartIcon />,
    title: 'Code Collaboration'
  },
  {
    description:
      'Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.',
    icon: <WalletIcon />,
    title: 'Project Management'
  },
  {
    description:
      'Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.',
    icon: <MagnifierIcon />,
    title: 'Task Automation'
  }
]

export const Services = () => {
  return (
    <section className='container py-24 sm:py-32'>
      <div className='grid place-items-center gap-8 lg:grid-cols-[1fr,1fr]'>
        <div>
          <h2 className='text-3xl font-bold md:text-4xl'>
            <span className='from-primary/60 to-primary bg-gradient-to-b bg-clip-text text-transparent'>
              Client-Centric{' '}
            </span>
            Services
          </h2>

          <p className='text-muted-foreground mt-4 mb-8 text-xl'>
            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Veritatis
            dolor.
          </p>

          <div className='flex flex-col gap-8'>
            {serviceList.map(({ description, icon, title }: ServiceProps) => (
              <Card key={title}>
                <CardHeader className='flex items-start justify-start gap-4 space-y-1 md:flex-row'>
                  <div className='bg-primary/20 mt-1 rounded-2xl p-1'>
                    {icon}
                  </div>
                  <div>
                    <CardTitle>{title}</CardTitle>
                    <CardDescription className='text-md mt-2'>
                      {description}
                    </CardDescription>
                  </div>
                </CardHeader>
              </Card>
            ))}
          </div>
        </div>

        <Image
          alt='About services'
          className='w-[300px] object-contain md:w-[500px] lg:w-[600px]'
          height={300}
          src={'/cube-leg.png'}
          width={500}
        />
      </div>
    </section>
  )
}
