import { ScrollButton } from '@archesai/ui/components/custom/scroll-button'

import { About } from '#app/(landing)/components/About'
import { Cta } from '#app/(landing)/components/Cta'
import { FAQ } from '#app/(landing)/components/FAQ'
import { Features } from '#app/(landing)/components/Features'
import { Footer } from '#app/(landing)/components/Footer'
import { Hero } from '#app/(landing)/components/Hero'
import { HowItWorks } from '#app/(landing)/components/HowItWorks'
import { Navbar } from '#app/(landing)/components/Navbar'
import { Newsletter } from '#app/(landing)/components/Newsletter'
import { Pricing } from '#app/(landing)/components/Pricing'
import { Services } from '#app/(landing)/components/Services'
import { Testimonials } from '#app/(landing)/components/Testimonials'

function LandingPage() {
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
      <Pricing />
      <Newsletter />
      <FAQ />
      <Footer />
      <ScrollButton />
    </>
  )
}

export default LandingPage
