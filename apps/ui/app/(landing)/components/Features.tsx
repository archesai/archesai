import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import Image from 'next/image'

interface FeatureProps {
  description: string
  image: string
  title: string
}

const features: FeatureProps[] = [
  {
    description:
      "Our platform ensures responsiveness across all devices, allowing you to access Arches AI's tools seamlessly from any screen size.",
    image: '/responsive-design.png', // Replace with relevant image showing responsive design
    title: 'Responsive Design'
  },
  {
    description:
      'Designed with users in mind, Arches AI offers an intuitive interface that simplifies complex workflows and maximizes productivity.',
    image: '/intuitive-ui.png', // Replace with relevant image depicting intuitive UI
    title: 'Intuitive User Interface'
  },
  {
    description:
      'Harness the power of artificial intelligence to gather insights from your data, helping you make smarter decisions faster.',
    image: '/ai-insights.png', // Replace with relevant image depicting AI-powered insights
    title: 'AI-Powered Insights'
  }
]

const featureList: string[] = [
  'Dark/Light Mode',
  'Advanced Analytics',
  'AI-Driven Automation',
  'Real-time Collaboration',
  'Customizable Dashboards',
  'Security and Compliance',
  'Responsive Design',
  'Seamless Integrations',
  'User Management'
]

export const Features = () => {
  return (
    <section
      className='container space-y-8 py-24 sm:py-32'
      id='features'
    >
      <h2 className='text-3xl font-bold md:text-center lg:text-4xl'>
        Many{' '}
        <span className='from-primary/60 to-primary bg-gradient-to-b bg-clip-text text-transparent'>
          Great Features
        </span>
      </h2>

      <div className='flex flex-wrap gap-4 md:justify-center'>
        {featureList.map((feature: string) => (
          <div key={feature}>
            <Badge
              className='text-sm'
              variant='secondary'
            >
              {feature}
            </Badge>
          </div>
        ))}
      </div>

      <div className='grid gap-8 md:grid-cols-2 lg:grid-cols-3'>
        {features.map(({ description, image, title }: FeatureProps) => (
          <Card key={title}>
            <CardHeader>
              <CardTitle>{title}</CardTitle>
            </CardHeader>

            <CardContent>{description}</CardContent>

            <CardFooter>
              <Image
                alt={title}
                className='mx-auto w-[200px] lg:w-[300px]'
                height={200}
                src={image}
                width={200}
              />
            </CardFooter>
          </Card>
        ))}
      </div>
    </section>
  )
}
