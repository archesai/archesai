import type { JSX } from "react"

import { useEffect, useState } from "react"
import { createFileRoute, Link } from "@tanstack/react-router"
import { motion } from "motion/react"

import { ArchesLogo } from "@archesai/ui/components/custom/arches-logo"
import {
  ArrowRightIcon,
  BarChartIcon,
  CheckCircle2Icon,
  ChevronRightIcon,
  LayersIcon,
  MenuIcon,
  ShieldIcon,
  StarIcon,
  UsersIcon,
  XCircleIcon,
  ZapIcon
} from "@archesai/ui/components/custom/icons"
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger
} from "@archesai/ui/components/shadcn/accordion"
import { Badge } from "@archesai/ui/components/shadcn/badge"
import { Button } from "@archesai/ui/components/shadcn/button"
import { Card, CardContent } from "@archesai/ui/components/shadcn/card"
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from "@archesai/ui/components/shadcn/tabs"

export const Route = createFileRoute("/landing/")({
  component: RouteComponent
})

export default function LandingPage(): JSX.Element {
  const [isScrolled, setIsScrolled] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  // const { resolvedTheme, setTheme } = useTheme()
  const [_mounted, setMounted] = useState(false)

  // const toggleTheme = useCallback(() => {
  //   setTheme(resolvedTheme === 'dark' ? 'light' : 'dark')
  // }, [resolvedTheme, setTheme])

  useEffect(() => {
    setMounted(true)
    const handleScroll = () => {
      if (window.scrollY > 10) {
        setIsScrolled(true)
      } else {
        setIsScrolled(false)
      }
    }

    window.addEventListener("scroll", handleScroll)
    return () => {
      window.removeEventListener("scroll", handleScroll)
    }
  }, [])

  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  }

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 }
  }

  const features = [
    {
      description:
        "Automate repetitive tasks and workflows to save time and reduce errors.",
      icon: <ZapIcon className="size-5" />,
      title: "Smart Automation"
    },
    {
      description:
        "Gain valuable insights with real-time data visualization and reporting.",
      icon: <BarChartIcon className="size-5" />,
      title: "Advanced Analytics"
    },
    {
      description:
        "Work together seamlessly with integrated communication tools.",
      icon: <UsersIcon className="size-5" />,
      title: "Team Collaboration"
    },
    {
      description:
        "Keep your data safe with end-to-end encryption and compliance features.",
      icon: <ShieldIcon className="size-5" />,
      title: "Enterprise Security"
    },
    {
      description:
        "Connect with your favorite tools through our extensive API ecosystem.",
      icon: <LayersIcon className="size-5" />,
      title: "Seamless Integration"
    },
    {
      description:
        "Get help whenever you need it with our dedicated support team.",
      icon: <StarIcon className="size-5" />,
      title: "24/7 Support"
    }
  ]

  return (
    <div className="flex min-h-[100dvh] flex-col bg-black">
      <header
        className={`sticky top-0 z-50 w-full backdrop-blur-lg transition-all duration-300 ${isScrolled ? "bg-background/80 shadow-sm" : "bg-transparent"}`}
      >
        <div className="container flex h-16 items-center justify-between">
          <ArchesLogo size="lg" />
          <nav className="hidden gap-8 md:flex">
            <Link
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              search={{ scrollTo: "features" }}
              to="/landing"
            >
              Features
            </Link>
            <Link
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              search={{ scrollTo: "testimonials" }}
              to="/landing"
            >
              Testimonials
            </Link>
            <Link
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              search={{ scrollTo: "pricing" }}
              to="/landing"
            >
              Pricing
            </Link>
            <Link
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              search={{ scrollTo: "faq" }}
              to="/landing"
            >
              FAQ
            </Link>
          </nav>
          <div className="hidden items-center gap-4 md:flex">
            {/* <Button
              className='rounded-full'
              onClick={toggleTheme}
              size='icon'
              variant='ghost'
            >
              {mounted && resolvedTheme === 'dark' ?
                <Sun className='size-[18px]' />
              : <Moon className='size-[18px]' />}
              <span className='sr-only'>Toggle theme</span>
            </Button> */}
            <Link
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
              to="/"
            >
              Log in
            </Link>
            <Button className="rounded-full">
              Get Started
              <ChevronRightIcon className="ml-1 size-4" />
            </Button>
          </div>
          <div className="flex items-center gap-4 md:hidden">
            {/* <Button
              className='rounded-full'
              onClick={toggleTheme}
              size='icon'
              variant='ghost'
            >
              {mounted && resolvedTheme === 'dark' ?
                <Sun className='size-[18px]' />
              : <Moon className='size-[18px]' />}
            </Button> */}
            <Button
              onClick={() => {
                setMobileMenuOpen(!mobileMenuOpen)
              }}
              size="icon"
              variant="ghost"
            >
              {mobileMenuOpen ? (
                <XCircleIcon className="size-5" />
              ) : (
                <MenuIcon className="size-5" />
              )}
              <span className="sr-only">Toggle menu</span>
            </Button>
          </div>
        </div>
        {/* Mobile menu */}
        {mobileMenuOpen && (
          <motion.div
            animate={{ opacity: 1, y: 0 }}
            className="absolute inset-x-0 top-16 border-b bg-background/95 backdrop-blur-lg md:hidden"
            exit={{ opacity: 0, y: -20 }}
            initial={{ opacity: 0, y: -20 }}
          >
            <div className="container flex flex-col gap-4 py-4">
              <Link
                className="py-2 text-sm font-medium"
                onClick={() => {
                  setMobileMenuOpen(false)
                }}
                search={{ scrollTo: "features" }}
                to="/landing"
              >
                Features
              </Link>
              <Link
                className="py-2 text-sm font-medium"
                onClick={() => {
                  setMobileMenuOpen(false)
                }}
                search={{ scrollTo: "testimonials" }}
                to="/landing"
              >
                Testimonials
              </Link>
              <Link
                className="py-2 text-sm font-medium"
                onClick={() => {
                  setMobileMenuOpen(false)
                }}
                search={{ scrollTo: "pricing" }}
                to="/landing"
              >
                Pricing
              </Link>
              <Link
                className="py-2 text-sm font-medium"
                onClick={() => {
                  setMobileMenuOpen(false)
                }}
                search={{ scrollTo: "faq" }}
                to="/landing"
              >
                FAQ
              </Link>
              <div className="flex flex-col gap-2 border-t pt-2">
                <Link
                  className="py-2 text-sm font-medium"
                  onClick={() => {
                    setMobileMenuOpen(false)
                  }}
                  to="/"
                >
                  Log in
                </Link>
                <Button className="rounded-full">
                  Get Started
                  <ChevronRightIcon className="ml-1 size-4" />
                </Button>
              </div>
            </div>
          </motion.div>
        )}
      </header>
      <main className="flex-1">
        {/* Hero Section */}
        <section className="w-full overflow-hidden py-20 md:py-32 lg:py-40">
          <div className="relative container px-4 md:px-6">
            <div className="absolute inset-0 -z-10 h-full w-full bg-white bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_110%)] bg-[size:4rem_4rem] dark:bg-black dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)]"></div>

            <motion.div
              animate={{ opacity: 1, y: 0 }}
              className="mx-auto mb-12 max-w-3xl text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
            >
              <Badge
                className="mb-4 rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                Launching Soon
              </Badge>
              <h1 className="mb-6 bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-4xl font-bold tracking-tight text-transparent md:text-5xl lg:text-6xl">
                Elevate Your Workflow with SaaSify
              </h1>
              <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground md:text-xl">
                The all-in-one platform that helps teams collaborate, automate,
                and deliver exceptional results. Streamline your processes and
                focus on what matters most.
              </p>
              <div className="flex flex-col justify-center gap-4 sm:flex-row">
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                >
                  Start Free Trial
                  <ArrowRightIcon className="ml-2 size-4" />
                </Button>
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                  variant="outline"
                >
                  Book a Demo
                </Button>
              </div>
              <div className="mt-6 flex items-center justify-center gap-4 text-sm text-muted-foreground">
                <div className="flex items-center gap-1">
                  <CheckCircle2Icon className="size-4 text-primary" />
                  <span>No credit card</span>
                </div>
                <div className="flex items-center gap-1">
                  <CheckCircle2Icon className="size-4 text-primary" />
                  <span>14-day trial</span>
                </div>
                <div className="flex items-center gap-1">
                  <CheckCircle2Icon className="size-4 text-primary" />
                  <span>Cancel anytime</span>
                </div>
              </div>
            </motion.div>

            <motion.div
              animate={{ opacity: 1, y: 0 }}
              className="relative mx-auto max-w-5xl"
              initial={{ opacity: 0, y: 40 }}
              transition={{ delay: 0.2, duration: 0.7 }}
            >
              <div className="overflow-hidden rounded-xl border border-border/40 bg-gradient-to-b from-background to-muted/20 shadow-2xl">
                <img
                  alt="SaaSify dashboard"
                  className="h-auto w-full"
                  height={720}
                  src="https://cdn.dribbble.com/userupload/12302729/file/original-fa372845e394ee85bebe0389b9d86871.png?resize=1504x1128&vertical=center"
                  width={1280}
                />
                <div className="absolute inset-0 rounded-xl ring-1 ring-black/10 ring-inset dark:ring-white/10"></div>
              </div>
              <div className="absolute -right-6 -bottom-6 -z-10 h-[300px] w-[300px] rounded-full bg-gradient-to-br from-primary/30 to-secondary/30 opacity-70 blur-3xl"></div>
              <div className="absolute -top-6 -left-6 -z-10 h-[300px] w-[300px] rounded-full bg-gradient-to-br from-secondary/30 to-primary/30 opacity-70 blur-3xl"></div>
            </motion.div>
          </div>
        </section>

        {/* Logos Section */}
        <section className="w-full border-y bg-muted/30 py-12">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <p className="text-sm font-medium text-muted-foreground">
                Trusted by innovative companies worldwide
              </p>
              <div className="flex flex-wrap items-center justify-center gap-8 md:gap-12 lg:gap-16">
                {[1, 2, 3, 4, 5].map((i) => (
                  <img
                    alt={`Company logo ${i.toString()}`}
                    className="h-8 w-auto opacity-70 grayscale transition-all hover:opacity-100 hover:grayscale-0"
                    height={60}
                    key={i}
                    src={"/placeholder-logo.svg"}
                    width={120}
                  />
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section
          className="w-full py-20 md:py-32"
          id="features"
        >
          <div className="container px-4 md:px-6">
            <motion.div
              className="mb-12 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                Features
              </Badge>
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl">
                Everything You Need to Succeed
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                Our comprehensive platform provides all the tools you need to
                streamline your workflow, boost productivity, and achieve your
                goals.
              </p>
            </motion.div>

            <motion.div
              className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3"
              initial="hidden"
              variants={container}
              viewport={{ once: true }}
              whileInView="show"
            >
              {features.map((feature, i) => (
                <motion.div
                  key={i}
                  variants={item}
                >
                  <Card className="h-full overflow-hidden border-border/40 bg-gradient-to-b from-background to-muted/10 backdrop-blur transition-all hover:shadow-md">
                    <CardContent className="flex h-full flex-col p-6">
                      <div className="mb-4 flex size-10 items-center justify-center rounded-full bg-primary/10 text-primary dark:bg-primary/20">
                        {feature.icon}
                      </div>
                      <h3 className="mb-2 text-xl font-bold">
                        {feature.title}
                      </h3>
                      <p className="text-muted-foreground">
                        {feature.description}
                      </p>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </motion.div>
          </div>
        </section>

        {/* How It Works Section */}
        <section className="relative w-full overflow-hidden bg-muted/30 py-20 md:py-32">
          <div className="absolute inset-0 -z-10 h-full w-full bg-white bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_50%,#000_40%,transparent_100%)] bg-[size:4rem_4rem] dark:bg-black dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)]"></div>

          <div className="relative container px-4 md:px-6">
            <motion.div
              className="mb-16 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                How It Works
              </Badge>
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl">
                Simple Process, Powerful Results
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                Get started in minutes and see the difference our platform can
                make for your business.
              </p>
            </motion.div>

            <div className="relative grid gap-8 md:grid-cols-3 md:gap-12">
              <div className="absolute top-1/2 right-0 left-0 z-0 hidden h-0.5 -translate-y-1/2 bg-gradient-to-r from-transparent via-border to-transparent md:block"></div>

              {[
                {
                  description:
                    "Sign up in seconds with just your email. No credit card required to get started.",
                  step: "01",
                  title: "Create Account"
                },
                {
                  description:
                    "Customize your workspace to match your team's unique workflow and requirements.",
                  step: "02",
                  title: "Configure Workspace"
                },
                {
                  description:
                    "Start using our powerful features to streamline processes and achieve your goals.",
                  step: "03",
                  title: "Boost Productivity"
                }
              ].map((step, i) => (
                <motion.div
                  className="relative z-10 flex flex-col items-center space-y-4 text-center"
                  initial={{ opacity: 0, y: 20 }}
                  key={i}
                  transition={{ delay: i * 0.1, duration: 0.5 }}
                  viewport={{ once: true }}
                  whileInView={{ opacity: 1, y: 0 }}
                >
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-gradient-to-br from-primary to-primary/70 text-xl font-bold text-primary-foreground shadow-lg">
                    {step.step}
                  </div>
                  <h3 className="text-xl font-bold">{step.title}</h3>
                  <p className="text-muted-foreground">{step.description}</p>
                </motion.div>
              ))}
            </div>
          </div>
        </section>

        {/* Testimonials Section */}
        <section
          className="w-full py-20 md:py-32"
          id="testimonials"
        >
          <div className="container px-4 md:px-6">
            <motion.div
              className="mb-12 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                Testimonials
              </Badge>
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl">
                Loved by Teams Worldwide
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                Don&apos;t just take our word for it. See what our customers
                have to say about their experience.
              </p>
            </motion.div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {[
                {
                  author: "Sarah Johnson",
                  quote:
                    "SaaSify has transformed how we manage our projects. The automation features have saved us countless hours of manual work.",
                  rating: 5,
                  role: "Project Manager, TechCorp"
                },
                {
                  author: "Michael Chen",
                  quote:
                    "The analytics dashboard provides insights we never had access to before. It's helped us make data-driven decisions that have improved our ROI.",
                  rating: 5,
                  role: "Marketing Director, GrowthLabs"
                },
                {
                  author: "Emily Rodriguez",
                  quote:
                    "Customer support is exceptional. Any time we've had an issue, the team has been quick to respond and resolve it. Couldn't ask for better service.",
                  rating: 5,
                  role: "Operations Lead, StartupX"
                },
                {
                  author: "David Kim",
                  quote:
                    "We've tried several similar solutions, but none compare to the ease of use and comprehensive features of SaaSify. It's been a game-changer.",
                  rating: 5,
                  role: "CEO, InnovateNow"
                },
                {
                  author: "Lisa Patel",
                  quote:
                    "The collaboration tools have made remote work so much easier for our team. We're more productive than ever despite being spread across different time zones.",
                  rating: 5,
                  role: "HR Director, RemoteFirst"
                },
                {
                  author: "James Wilson",
                  quote:
                    "Implementation was seamless, and the ROI was almost immediate. We've reduced our operational costs by 30% since switching to SaaSify.",
                  rating: 5,
                  role: "COO, ScaleUp Inc"
                }
              ].map((testimonial, i) => (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  key={i}
                  transition={{ delay: i * 0.05, duration: 0.5 }}
                  viewport={{ once: true }}
                  whileInView={{ opacity: 1, y: 0 }}
                >
                  <Card className="h-full overflow-hidden border-border/40 bg-gradient-to-b from-background to-muted/10 backdrop-blur transition-all hover:shadow-md">
                    <CardContent className="flex h-full flex-col p-6">
                      <div className="mb-4 flex">
                        {Array(testimonial.rating)
                          .fill(0)
                          .map((_, j) => (
                            <StarIcon
                              className="size-4 fill-yellow-500 text-yellow-500"
                              key={j}
                            />
                          ))}
                      </div>
                      <p className="mb-6 flex-grow text-lg">
                        {testimonial.quote}
                      </p>
                      <div className="mt-auto flex items-center gap-4 border-t border-border/40 pt-4">
                        <div className="flex size-10 items-center justify-center rounded-full bg-muted font-medium text-foreground">
                          {testimonial.author.charAt(0)}
                        </div>
                        <div>
                          <p className="font-medium">{testimonial.author}</p>
                          <p className="text-sm text-muted-foreground">
                            {testimonial.role}
                          </p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          </div>
        </section>

        {/* Pricing Section */}
        <section
          className="relative w-full overflow-hidden bg-muted/30 py-20 md:py-32"
          id="pricing"
        >
          <div className="absolute inset-0 -z-10 h-full w-full bg-white bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_50%,#000_40%,transparent_100%)] bg-[size:4rem_4rem] dark:bg-black dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)]"></div>

          <div className="relative container px-4 md:px-6">
            <motion.div
              className="mb-12 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                Pricing
              </Badge>
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl">
                Simple, Transparent Pricing
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                Choose the plan that&apos;s right for your business. All plans
                include a 14-day free trial.
              </p>
            </motion.div>

            <div className="mx-auto max-w-5xl">
              <Tabs
                className="w-full"
                defaultValue="monthly"
              >
                <div className="mb-8 flex justify-center">
                  <TabsList className="rounded-full p-1">
                    <TabsTrigger
                      className="rounded-full px-6"
                      value="monthly"
                    >
                      Monthly
                    </TabsTrigger>
                    <TabsTrigger
                      className="rounded-full px-6"
                      value="annually"
                    >
                      Annually (Save 20%)
                    </TabsTrigger>
                  </TabsList>
                </div>
                <TabsContent value="monthly">
                  <div className="grid gap-6 lg:grid-cols-3 lg:gap-8">
                    {[
                      {
                        cta: "Start Free Trial",
                        description: "Perfect for small teams and startups.",
                        features: [
                          "Up to 5 team members",
                          "Basic analytics",
                          "5GB storage",
                          "Email support"
                        ],
                        name: "Starter",
                        price: "$29"
                      },
                      {
                        cta: "Start Free Trial",
                        description: "Ideal for growing businesses.",
                        features: [
                          "Up to 20 team members",
                          "Advanced analytics",
                          "25GB storage",
                          "Priority email support",
                          "API access"
                        ],
                        name: "Professional",
                        popular: true,
                        price: "$79"
                      },
                      {
                        cta: "Contact Sales",
                        description:
                          "For large organizations with complex needs.",
                        features: [
                          "Unlimited team members",
                          "Custom analytics",
                          "Unlimited storage",
                          "24/7 phone & email support",
                          "Advanced API access",
                          "Custom integrations"
                        ],
                        name: "Enterprise",
                        price: "$199"
                      }
                    ].map((plan, i) => (
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        key={i}
                        transition={{ delay: i * 0.1, duration: 0.5 }}
                        viewport={{ once: true }}
                        whileInView={{ opacity: 1, y: 0 }}
                      >
                        <Card
                          className={`relative h-full overflow-hidden ${plan.popular ? "border-primary shadow-lg" : "border-border/40 shadow-md"} bg-gradient-to-b from-background to-muted/10 backdrop-blur`}
                        >
                          {plan.popular && (
                            <div className="absolute top-0 right-0 rounded-bl-lg bg-primary px-3 py-1 text-xs font-medium text-primary-foreground">
                              Most Popular
                            </div>
                          )}
                          <CardContent className="flex h-full flex-col p-6">
                            <h3 className="text-2xl font-bold">{plan.name}</h3>
                            <div className="mt-4 flex items-baseline">
                              <span className="text-4xl font-bold">
                                {plan.price}
                              </span>
                              <span className="ml-1 text-muted-foreground">
                                /month
                              </span>
                            </div>
                            <p className="mt-2 text-muted-foreground">
                              {plan.description}
                            </p>
                            <ul className="my-6 flex-grow space-y-3">
                              {plan.features.map((feature, j) => (
                                <li
                                  className="flex items-center"
                                  key={j}
                                >
                                  <CheckCircle2Icon className="mr-2 size-4 text-primary" />
                                  <span>{feature}</span>
                                </li>
                              ))}
                            </ul>
                            <Button
                              className={`mt-auto w-full rounded-full ${plan.popular ? "bg-primary hover:bg-primary/90" : "bg-muted hover:bg-muted/80"}`}
                              variant={plan.popular ? "default" : "outline"}
                            >
                              {plan.cta}
                            </Button>
                          </CardContent>
                        </Card>
                      </motion.div>
                    ))}
                  </div>
                </TabsContent>
                <TabsContent value="annually">
                  <div className="grid gap-6 lg:grid-cols-3 lg:gap-8">
                    {[
                      {
                        cta: "Start Free Trial",
                        description: "Perfect for small teams and startups.",
                        features: [
                          "Up to 5 team members",
                          "Basic analytics",
                          "5GB storage",
                          "Email support"
                        ],
                        name: "Starter",
                        price: "$23"
                      },
                      {
                        cta: "Start Free Trial",
                        description: "Ideal for growing businesses.",
                        features: [
                          "Up to 20 team members",
                          "Advanced analytics",
                          "25GB storage",
                          "Priority email support",
                          "API access"
                        ],
                        name: "Professional",
                        popular: true,
                        price: "$63"
                      },
                      {
                        cta: "Contact Sales",
                        description:
                          "For large organizations with complex needs.",
                        features: [
                          "Unlimited team members",
                          "Custom analytics",
                          "Unlimited storage",
                          "24/7 phone & email support",
                          "Advanced API access",
                          "Custom integrations"
                        ],
                        name: "Enterprise",
                        price: "$159"
                      }
                    ].map((plan, i) => (
                      <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        key={i}
                        transition={{ delay: i * 0.1, duration: 0.5 }}
                        viewport={{ once: true }}
                        whileInView={{ opacity: 1, y: 0 }}
                      >
                        <Card
                          className={`relative h-full overflow-hidden ${plan.popular ? "border-primary shadow-lg" : "border-border/40 shadow-md"} bg-gradient-to-b from-background to-muted/10 backdrop-blur`}
                        >
                          {plan.popular && (
                            <div className="absolute top-0 right-0 rounded-bl-lg bg-primary px-3 py-1 text-xs font-medium text-primary-foreground">
                              Most Popular
                            </div>
                          )}
                          <CardContent className="flex h-full flex-col p-6">
                            <h3 className="text-2xl font-bold">{plan.name}</h3>
                            <div className="mt-4 flex items-baseline">
                              <span className="text-4xl font-bold">
                                {plan.price}
                              </span>
                              <span className="ml-1 text-muted-foreground">
                                /month
                              </span>
                            </div>
                            <p className="mt-2 text-muted-foreground">
                              {plan.description}
                            </p>
                            <ul className="my-6 flex-grow space-y-3">
                              {plan.features.map((feature, j) => (
                                <li
                                  className="flex items-center"
                                  key={j}
                                >
                                  <CheckCircle2Icon className="mr-2 size-4 text-primary" />
                                  <span>{feature}</span>
                                </li>
                              ))}
                            </ul>
                            <Button
                              className={`mt-auto w-full rounded-full ${plan.popular ? "bg-primary hover:bg-primary/90" : "bg-muted hover:bg-muted/80"}`}
                              variant={plan.popular ? "default" : "outline"}
                            >
                              {plan.cta}
                            </Button>
                          </CardContent>
                        </Card>
                      </motion.div>
                    ))}
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          </div>
        </section>

        {/* FAQ Section */}
        <section
          className="w-full py-20 md:py-32"
          id="faq"
        >
          <div className="container px-4 md:px-6">
            <motion.div
              className="mb-12 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 text-sm font-medium"
                variant="secondary"
              >
                FAQ
              </Badge>
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl">
                Frequently Asked Questions
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                Find answers to common questions about our platform.
              </p>
            </motion.div>

            <div className="mx-auto max-w-3xl">
              <Accordion
                className="w-full"
                collapsible
                type="single"
              >
                {[
                  {
                    answer:
                      "Our 14-day free trial gives you full access to all features of your selected plan. No credit card is required to sign up, and you can cancel at any time during the trial period with no obligation.",
                    question: "How does the 14-day free trial work?"
                  },
                  {
                    answer:
                      "Yes, you can upgrade or downgrade your plan at any time. If you upgrade, the new pricing will be prorated for the remainder of your billing cycle. If you downgrade, the new pricing will take effect at the start of your next billing cycle.",
                    question: "Can I change plans later?"
                  },
                  {
                    answer:
                      "The number of usersIcon depends on your plan. The Starter plan allows up to 5 team members, the Professional plan allows up to 20, and the Enterprise plan has no limit on team members.",
                    question:
                      "Is there a limit to how many usersIcon I can add?"
                  },
                  {
                    answer:
                      "Yes, we offer special pricing for nonprofits, educational institutions, and open-source projects. Please contact our sales team for more information.",
                    question:
                      "Do you offer discounts for nonprofits or educational institutions?"
                  },
                  {
                    answer:
                      "We take security very seriously. All data is encrypted both in transit and at rest. We use industry-standard security practices and regularly undergo security audits. Our platform is compliant with GDPR, CCPA, and other relevant regulations.",
                    question: "How secure is my data?"
                  },
                  {
                    answer:
                      "Support varies by plan. All plans include email support, with the Professional plan offering priority email support. The Enterprise plan includes 24/7 phone and email support. We also have an extensive knowledge base and community forum available to all users.",
                    question: "What kind of support do you offer?"
                  }
                ].map((faq, i) => (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    key={i}
                    transition={{ delay: i * 0.05, duration: 0.3 }}
                    viewport={{ once: true }}
                    whileInView={{ opacity: 1, y: 0 }}
                  >
                    <AccordionItem
                      className="border-b border-border/40 py-2"
                      value={`item-${i.toString()}`}
                    >
                      <AccordionTrigger className="text-left font-medium hover:no-underline">
                        {faq.question}
                      </AccordionTrigger>
                      <AccordionContent className="text-muted-foreground">
                        {faq.answer}
                      </AccordionContent>
                    </AccordionItem>
                  </motion.div>
                ))}
              </Accordion>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="relative w-full overflow-hidden bg-gradient-to-br from-primary to-primary/80 py-20 text-primary-foreground md:py-32">
          <div className="absolute inset-0 -z-10 bg-[linear-gradient(to_right,#ffffff10_1px,transparent_1px),linear-gradient(to_bottom,#ffffff10_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
          <div className="absolute -top-24 -left-24 h-64 w-64 rounded-full bg-white/10 blur-3xl"></div>
          <div className="absolute -right-24 -bottom-24 h-64 w-64 rounded-full bg-white/10 blur-3xl"></div>

          <div className="relative container px-4 md:px-6">
            <motion.div
              className="flex flex-col items-center justify-center space-y-6 text-center"
              initial={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.5 }}
              viewport={{ once: true }}
              whileInView={{ opacity: 1, y: 0 }}
            >
              <h2 className="text-3xl font-bold tracking-tight md:text-4xl lg:text-5xl">
                Ready to Transform Your Workflow?
              </h2>
              <p className="mx-auto max-w-[700px] text-primary-foreground/80 md:text-xl">
                Join thousands of satisfied customers who have streamlined their
                processes and boosted productivity with our platform.
              </p>
              <div className="mt-4 flex flex-col gap-4 sm:flex-row">
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                  variant="secondary"
                >
                  Start Free Trial
                  <ArrowRightIcon className="ml-2 size-4" />
                </Button>
                <Button
                  className="h-12 rounded-full border-white bg-transparent px-8 text-base text-white hover:bg-white/10"
                  size="lg"
                  variant="outline"
                >
                  Schedule a Demo
                </Button>
              </div>
              <p className="mt-4 text-sm text-primary-foreground/80">
                No credit card required. 14-day free trial. Cancel anytime.
              </p>
            </motion.div>
          </div>
        </section>
      </main>
      <footer className="w-full border-t bg-background/95 backdrop-blur-sm">
        <div className="container flex flex-col gap-8 px-4 py-10 md:px-6 lg:py-16">
          <div className="grid gap-8 sm:grid-cols-2 md:grid-cols-4">
            <div className="space-y-4">
              <div className="flex items-center gap-2 font-bold">
                <div className="flex size-8 items-center justify-center rounded-lg bg-gradient-to-br from-primary to-primary/70 text-primary-foreground">
                  S
                </div>
                <span>SaaSify</span>
              </div>
              <p className="text-sm text-muted-foreground">
                Streamline your workflow with our all-in-one SaaS platform.
                Boost productivity and scale your business.
              </p>
              <div className="flex gap-4">
                <Link
                  className="text-muted-foreground transition-colors hover:text-foreground"
                  to="/landing"
                >
                  <svg
                    className="size-5"
                    fill="none"
                    height="24"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    viewBox="0 0 24 24"
                    width="24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path>
                  </svg>
                  <span className="sr-only">Facebook</span>
                </Link>
                <Link
                  className="text-muted-foreground transition-colors hover:text-foreground"
                  to="/landing"
                >
                  <svg
                    className="size-5"
                    fill="none"
                    height="24"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    viewBox="0 0 24 24"
                    width="24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"></path>
                  </svg>
                  <span className="sr-only">Twitter</span>
                </Link>
                <Link
                  className="text-muted-foreground transition-colors hover:text-foreground"
                  to="/landing"
                >
                  <svg
                    className="size-5"
                    fill="none"
                    height="24"
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    viewBox="0 0 24 24"
                    width="24"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path>
                    <rect
                      height="12"
                      width="4"
                      x="2"
                      y="9"
                    ></rect>
                    <circle
                      cx="4"
                      cy="4"
                      r="2"
                    ></circle>
                  </svg>
                  <span className="sr-only">LinkedIn</span>
                </Link>
              </div>
            </div>
            <div className="space-y-4">
              <h4 className="text-sm font-bold">Product</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    search={{ scrollTo: "features" }}
                    to="/landing"
                  >
                    Features
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    search={{ scrollTo: "pricing" }}
                    to="/landing"
                  >
                    Pricing
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Integrations
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    API
                  </Link>
                </li>
              </ul>
            </div>
            <div className="space-y-4">
              <h4 className="text-sm font-bold">Resources</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Documentation
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Guides
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Blog
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Support
                  </Link>
                </li>
              </ul>
            </div>
            <div className="space-y-4">
              <h4 className="text-sm font-bold">Company</h4>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    About
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Careers
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Privacy Policy
                  </Link>
                </li>
                <li>
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    to="/"
                  >
                    Terms of Service
                  </Link>
                </li>
              </ul>
            </div>
          </div>
          <div className="flex flex-col items-center justify-between gap-4 border-t border-border/40 pt-8 sm:flex-row">
            <p className="text-xs text-muted-foreground">
              &copy; {new Date().getFullYear()} SaaSify. All rights reserved.
            </p>
            <div className="flex gap-4">
              <Link
                className="text-xs text-muted-foreground transition-colors hover:text-foreground"
                to="/"
              >
                Privacy Policy
              </Link>
              <Link
                className="text-xs text-muted-foreground transition-colors hover:text-foreground"
                to="/"
              >
                Terms of Service
              </Link>
              <Link
                className="text-xs text-muted-foreground transition-colors hover:text-foreground"
                to="/"
              >
                Cookie Policy
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}

function RouteComponent() {
  return <LandingPage />
}
