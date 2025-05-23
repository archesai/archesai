import { ScrollButton } from '@archesai/ui/components/custom/scroll-button'

import { About } from '#components/About'
import { Cta } from '#components/Cta'
import { FAQ } from '#components/FAQ'
import { Features } from '#components/Features'
import { Footer } from '#components/Footer'
import { Hero } from '#components/Hero'
import { HowItWorks } from '#components/HowItWorks'
import { Navbar } from '#components/Navbar'
import { Newsletter } from '#components/Newsletter'
// import { Pricing } from '#components/Pricing'
import { Services } from '#components/Services'
import { Testimonials } from '#components/Testimonials'

export default function LandingPage() {
  return (
    <>
      <Navbar />
      <Hero />
      <About />
      <HowItWorks />
      <Features />
      <Services />
      <Cta />
      <Testimonials />
      {/* <Pricing /> FIXME */}
      <Newsletter />
      <FAQ />
      <Footer />
      <ScrollButton />
    </>
  )
}
