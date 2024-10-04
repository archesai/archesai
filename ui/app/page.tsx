"use client";
import { About } from "../components/landing/About";
import { Cta } from "../components/landing/Cta";
import { FAQ } from "../components/landing/FAQ";
import { Features } from "../components/landing/Features";
import { Footer } from "../components/landing/Footer";
import { Hero } from "../components/landing/Hero";
import { HowItWorks } from "../components/landing/HowItWorks";
import { Navbar } from "../components/landing/Navbar";
import { Newsletter } from "../components/landing/Newsletter";
import { Pricing } from "../components/landing/Pricing";
import { ScrollToTop } from "../components/landing/ScrollToTop";
import { Services } from "../components/landing/Services";
import { Testimonials } from "../components/landing/Testimonials";

function App() {
  return (
    <>
      <Navbar />
      {/* <div className="min-h-screen relative w-screen -mt-14 ">
        <Image
          alt="Arches AI - AI-Driven Solutions"
          className=""
          layout="fill"
          objectFit="cover"
          priority={true}
          src="/landing-compressed-2.jpg"
        />
        <div className="absolute bg-black bg-opacity-70 bottom-0 left-0 right-0 top-0" />
      </div> */}
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
      <ScrollToTop />
    </>
  );
}

export default App;
