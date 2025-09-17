import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
  ArchesLogo,
  ArrowRightIcon,
  Badge,
  BarChartIcon,
  Button,
  Card,
  CardContent,
  CheckCircle2Icon,
  ChevronRightIcon,
  LayersIcon,
  MenuIcon,
  ShieldIcon,
  StarIcon,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
  UsersIcon,
  XCircleIcon,
  ZapIcon,
} from "@archesai/ui";
import { motion } from "motion/react";
import type { JSX } from "react";
import { useEffect, useState } from "react";
import { Link } from "zudoku/router";
import type { LandingContent } from "./landing_content";
import { defaultContent } from "./landing_content";

// Icon mapping helper
const getIcon = (iconName: string) => {
  const icons: Record<string, JSX.Element> = {
    BarChartIcon: <BarChartIcon className="size-5" />,
    LayersIcon: <LayersIcon className="size-5" />,
    ShieldIcon: <ShieldIcon className="size-5" />,
    StarIcon: <StarIcon className="size-5" />,
    UsersIcon: <UsersIcon className="size-5" />,
    ZapIcon: <ZapIcon className="size-5" />,
  };
  return icons[iconName] || <StarIcon className="size-5" />;
};

interface LandingPageProps {
  content?: LandingContent;
}

export function LandingPage({
  content = defaultContent,
}: LandingPageProps): JSX.Element {
  const [isScrolled, setIsScrolled] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [_mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    const handleScroll = () => {
      if (window.scrollY > 10) {
        setIsScrolled(true);
      } else {
        setIsScrolled(false);
      }
    };

    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  const container = {
    hidden: {
      opacity: 0,
    },
    show: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const item = {
    hidden: {
      opacity: 0,
      y: 20,
    },
    show: {
      opacity: 1,
      y: 0,
    },
  };

  return (
    <div className="flex min-h-[100dvh] flex-col bg-black">
      <header
        className={`sticky top-0 z-50 w-full backdrop-blur-lg transition-all duration-300 ${isScrolled ? "bg-background/80 shadow-sm" : "bg-transparent"}`}
      >
        <div className="container flex h-16 items-center justify-between">
          <ArchesLogo />
          <nav className="hidden gap-8 md:flex">
            {content.navigation.links.map((link) => (
              <Link
                className="font-medium text-muted-foreground text-sm transition-colors hover:text-foreground"
                key={link.label}
                {...(link.scrollTo
                  ? {
                      search: {
                        scrollTo: link.scrollTo,
                      },
                    }
                  : {})}
                to={link.to}
              >
                {link.label}
              </Link>
            ))}
          </nav>
          <div className="hidden items-center gap-4 md:flex">
            <Link
              className="font-medium text-muted-foreground text-sm transition-colors hover:text-foreground"
              to={content.navigation.buttons.login.to}
            >
              {content.navigation.buttons.login.label}
            </Link>
            <Button className="rounded-full">
              {content.navigation.buttons.getStarted.label}
              <ChevronRightIcon className="ml-1 size-4" />
            </Button>
          </div>
          <div className="flex items-center gap-4 md:hidden">
            <Button
              onClick={() => {
                setMobileMenuOpen(!mobileMenuOpen);
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
            animate={{
              opacity: 1,
              y: 0,
            }}
            className="absolute inset-x-0 top-16 border-b bg-background/95 backdrop-blur-lg md:hidden"
            exit={{
              opacity: 0,
              y: -20,
            }}
            initial={{
              opacity: 0,
              y: -20,
            }}
          >
            <div className="container flex flex-col gap-4 py-4">
              {content.navigation.links.map((link) => (
                <Link
                  className="py-2 font-medium text-sm"
                  key={link.label}
                  onClick={() => {
                    setMobileMenuOpen(false);
                  }}
                  {...(link.scrollTo
                    ? {
                        search: {
                          scrollTo: link.scrollTo,
                        },
                      }
                    : {})}
                  to={link.to}
                >
                  {link.label}
                </Link>
              ))}
              <div className="flex flex-col gap-2 border-t pt-2">
                <Link
                  className="py-2 font-medium text-sm"
                  onClick={() => {
                    setMobileMenuOpen(false);
                  }}
                  to={content.navigation.buttons.login.to}
                >
                  {content.navigation.buttons.login.label}
                </Link>
                <Button className="rounded-full">
                  {content.navigation.buttons.getStarted.label}
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
          <div className="container relative px-4 md:px-6">
            <div className="-z-10 absolute inset-0 h-full w-full bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] bg-[size:4rem_4rem] bg-white [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_110%)] dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)] dark:bg-black"></div>

            <motion.div
              animate={{
                opacity: 1,
                y: 0,
              }}
              className="mx-auto mb-12 max-w-3xl text-center"
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
            >
              <Badge
                className="mb-4 rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.hero.badge}
              </Badge>
              <h1 className="mb-6 bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text font-bold text-4xl text-transparent tracking-tight md:text-5xl lg:text-6xl">
                {content.hero.title}
              </h1>
              <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground md:text-xl">
                {content.hero.subtitle}
              </p>
              <div className="flex flex-col justify-center gap-4 sm:flex-row">
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                >
                  {content.hero.buttons.primary.label}
                  <ArrowRightIcon className="ml-2 size-4" />
                </Button>
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                  variant="outline"
                >
                  {content.hero.buttons.secondary.label}
                </Button>
              </div>
              <div className="mt-6 flex items-center justify-center gap-4 text-muted-foreground text-sm">
                {content.hero.benefits.map((benefit) => (
                  <div
                    className="flex items-center gap-1"
                    key={benefit}
                  >
                    <CheckCircle2Icon className="size-4 text-primary" />
                    <span>{benefit}</span>
                  </div>
                ))}
              </div>
            </motion.div>

            <motion.div
              animate={{
                opacity: 1,
                y: 0,
              }}
              className="relative mx-auto max-w-5xl"
              initial={{
                opacity: 0,
                y: 40,
              }}
              transition={{
                delay: 0.2,
                duration: 0.7,
              }}
            >
              <div className="overflow-hidden rounded-xl border border-border/40 bg-gradient-to-b from-background to-muted/20 shadow-2xl">
                <img
                  alt={content.hero.image.alt}
                  className="h-auto w-full"
                  height={720}
                  src={content.hero.image.src}
                  width={1280}
                />
                <div className="absolute inset-0 rounded-xl ring-1 ring-black/10 ring-inset dark:ring-white/10"></div>
              </div>
              <div className="-right-6 -bottom-6 -z-10 absolute h-[300px] w-[300px] rounded-full bg-gradient-to-br from-primary/30 to-secondary/30 opacity-70 blur-3xl"></div>
              <div className="-top-6 -left-6 -z-10 absolute h-[300px] w-[300px] rounded-full bg-gradient-to-br from-secondary/30 to-primary/30 opacity-70 blur-3xl"></div>
            </motion.div>
          </div>
        </section>

        {/* Logos Section */}
        <section className="w-full border-y bg-muted/30 py-12">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center justify-center space-y-4 text-center">
              <p className="font-medium text-muted-foreground text-sm">
                {content.logos.title}
              </p>
              <div className="flex flex-wrap items-center justify-center gap-8 md:gap-12 lg:gap-16">
                {content.logos.logos.map((logo) => (
                  <img
                    alt={logo.alt}
                    className="h-8 w-auto opacity-70 grayscale transition-all hover:opacity-100 hover:grayscale-0"
                    height={60}
                    key={logo.alt}
                    src={logo.src}
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
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.features.badge}
              </Badge>
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl">
                {content.features.title}
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                {content.features.subtitle}
              </p>
            </motion.div>

            <motion.div
              className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3"
              initial="hidden"
              variants={container}
              viewport={{
                once: true,
              }}
              whileInView="show"
            >
              {content.features.features.map((feature) => (
                <motion.div
                  key={feature.title}
                  variants={item}
                >
                  <Card className="h-full overflow-hidden border-border/40 bg-gradient-to-b from-background to-muted/10 backdrop-blur transition-all hover:shadow-md">
                    <CardContent className="flex h-full flex-col p-6">
                      <div className="mb-4 flex size-10 items-center justify-center rounded-full bg-primary/10 text-primary dark:bg-primary/20">
                        {getIcon(feature.icon)}
                      </div>
                      <h3 className="mb-2 font-bold text-xl">
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
          <div className="-z-10 absolute inset-0 h-full w-full bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] bg-[size:4rem_4rem] bg-white [mask-image:radial-gradient(ellipse_80%_50%_at_50%_50%,#000_40%,transparent_100%)] dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)] dark:bg-black"></div>

          <div className="container relative px-4 md:px-6">
            <motion.div
              className="mb-16 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.howItWorks.badge}
              </Badge>
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl">
                {content.howItWorks.title}
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                {content.howItWorks.subtitle}
              </p>
            </motion.div>

            <div className="relative grid gap-8 md:grid-cols-3 md:gap-12">
              <div className="-translate-y-1/2 absolute top-1/2 right-0 left-0 z-0 hidden h-0.5 bg-gradient-to-r from-transparent via-border to-transparent md:block"></div>

              {content.howItWorks.steps.map((step, i) => (
                <motion.div
                  className="relative z-10 flex flex-col items-center space-y-4 text-center"
                  initial={{
                    opacity: 0,
                    y: 20,
                  }}
                  key={step.step}
                  transition={{
                    delay: i * 0.1,
                    duration: 0.5,
                  }}
                  viewport={{
                    once: true,
                  }}
                  whileInView={{
                    opacity: 1,
                    y: 0,
                  }}
                >
                  <div className="flex h-16 w-16 items-center justify-center rounded-full bg-gradient-to-br from-primary to-primary/70 font-bold text-primary-foreground text-xl shadow-lg">
                    {step.step}
                  </div>
                  <h3 className="font-bold text-xl">{step.title}</h3>
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
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.testimonials.badge}
              </Badge>
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl">
                {content.testimonials.title}
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                {content.testimonials.subtitle}
              </p>
            </motion.div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {content.testimonials.testimonials.map((testimonial, i) => (
                <motion.div
                  initial={{
                    opacity: 0,
                    y: 20,
                  }}
                  key={testimonial.author}
                  transition={{
                    delay: i * 0.05,
                    duration: 0.5,
                  }}
                  viewport={{
                    once: true,
                  }}
                  whileInView={{
                    opacity: 1,
                    y: 0,
                  }}
                >
                  <Card className="h-full overflow-hidden border-border/40 bg-gradient-to-b from-background to-muted/10 backdrop-blur transition-all hover:shadow-md">
                    <CardContent className="flex h-full flex-col p-6">
                      <div className="mb-4 flex">
                        {Array(testimonial.rating)
                          .fill(0)
                          .map((_, j) => (
                            <StarIcon
                              className="size-4 fill-yellow-500 text-yellow-500"
                              key={`star-${testimonial.author}-${j}`}
                            />
                          ))}
                      </div>
                      <p className="mb-6 flex-grow text-lg">
                        {testimonial.quote}
                      </p>
                      <div className="mt-auto flex items-center gap-4 border-border/40 border-t pt-4">
                        <div className="flex size-10 items-center justify-center rounded-full bg-muted font-medium text-foreground">
                          {testimonial.author.charAt(0)}
                        </div>
                        <div>
                          <p className="font-medium">{testimonial.author}</p>
                          <p className="text-muted-foreground text-sm">
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
          <div className="-z-10 absolute inset-0 h-full w-full bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] bg-[size:4rem_4rem] bg-white [mask-image:radial-gradient(ellipse_80%_50%_at_50%_50%,#000_40%,transparent_100%)] dark:bg-[linear-gradient(to_right,#1f1f1f_1px,transparent_1px),linear-gradient(to_bottom,#1f1f1f_1px,transparent_1px)] dark:bg-black"></div>

          <div className="container relative px-4 md:px-6">
            <motion.div
              className="mb-12 flex flex-col items-center justify-center space-y-4 text-center"
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.pricing.badge}
              </Badge>
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl">
                {content.pricing.title}
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                {content.pricing.subtitle}
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
                      {content.pricing.tabs.monthly.label}
                    </TabsTrigger>
                    <TabsTrigger
                      className="rounded-full px-6"
                      value="annually"
                    >
                      {content.pricing.tabs.annually.label} (
                      {content.pricing.tabs.annually.savingsText})
                    </TabsTrigger>
                  </TabsList>
                </div>
                <TabsContent value="monthly">
                  <div className="grid gap-6 lg:grid-cols-3 lg:gap-8">
                    {content.pricing.tabs.monthly.plans.map((plan, i) => (
                      <motion.div
                        initial={{
                          opacity: 0,
                          y: 20,
                        }}
                        key={plan.name}
                        transition={{
                          delay: i * 0.1,
                          duration: 0.5,
                        }}
                        viewport={{
                          once: true,
                        }}
                        whileInView={{
                          opacity: 1,
                          y: 0,
                        }}
                      >
                        <Card
                          className={`relative h-full overflow-hidden ${plan.popular ? "border-primary shadow-lg" : "border-border/40 shadow-md"} bg-gradient-to-b from-background to-muted/10 backdrop-blur`}
                        >
                          {plan.popular && (
                            <div className="absolute top-0 right-0 rounded-bl-lg bg-primary px-3 py-1 font-medium text-primary-foreground text-xs">
                              Most Popular
                            </div>
                          )}
                          <CardContent className="flex h-full flex-col p-6">
                            <h3 className="font-bold text-2xl">{plan.name}</h3>
                            <div className="mt-4 flex items-baseline">
                              <span className="font-bold text-4xl">
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
                              {plan.features.map((feature) => (
                                <li
                                  className="flex items-center"
                                  key={feature}
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
                    {content.pricing.tabs.annually.plans.map((plan, i) => (
                      <motion.div
                        initial={{
                          opacity: 0,
                          y: 20,
                        }}
                        key={plan.name}
                        transition={{
                          delay: i * 0.1,
                          duration: 0.5,
                        }}
                        viewport={{
                          once: true,
                        }}
                        whileInView={{
                          opacity: 1,
                          y: 0,
                        }}
                      >
                        <Card
                          className={`relative h-full overflow-hidden ${plan.popular ? "border-primary shadow-lg" : "border-border/40 shadow-md"} bg-gradient-to-b from-background to-muted/10 backdrop-blur`}
                        >
                          {plan.popular && (
                            <div className="absolute top-0 right-0 rounded-bl-lg bg-primary px-3 py-1 font-medium text-primary-foreground text-xs">
                              Most Popular
                            </div>
                          )}
                          <CardContent className="flex h-full flex-col p-6">
                            <h3 className="font-bold text-2xl">{plan.name}</h3>
                            <div className="mt-4 flex items-baseline">
                              <span className="font-bold text-4xl">
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
                              {plan.features.map((feature) => (
                                <li
                                  className="flex items-center"
                                  key={feature}
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
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <Badge
                className="rounded-full px-4 py-1.5 font-medium text-sm"
                variant="secondary"
              >
                {content.faq.badge}
              </Badge>
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl">
                {content.faq.title}
              </h2>
              <p className="max-w-[800px] text-muted-foreground md:text-lg">
                {content.faq.subtitle}
              </p>
            </motion.div>

            <div className="mx-auto max-w-3xl">
              <Accordion
                className="w-full"
                collapsible
                type="single"
              >
                {content.faq.questions.map((faq, i) => (
                  <motion.div
                    initial={{
                      opacity: 0,
                      y: 10,
                    }}
                    key={faq.question}
                    transition={{
                      delay: i * 0.05,
                      duration: 0.3,
                    }}
                    viewport={{
                      once: true,
                    }}
                    whileInView={{
                      opacity: 1,
                      y: 0,
                    }}
                  >
                    <AccordionItem
                      className="border-border/40 border-b py-2"
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
          <div className="-z-10 absolute inset-0 bg-[linear-gradient(to_right,#ffffff10_1px,transparent_1px),linear-gradient(to_bottom,#ffffff10_1px,transparent_1px)] bg-[size:4rem_4rem]"></div>
          <div className="-top-24 -left-24 absolute h-64 w-64 rounded-full bg-white/10 blur-3xl"></div>
          <div className="-right-24 -bottom-24 absolute h-64 w-64 rounded-full bg-white/10 blur-3xl"></div>

          <div className="container relative px-4 md:px-6">
            <motion.div
              className="flex flex-col items-center justify-center space-y-6 text-center"
              initial={{
                opacity: 0,
                y: 20,
              }}
              transition={{
                duration: 0.5,
              }}
              viewport={{
                once: true,
              }}
              whileInView={{
                opacity: 1,
                y: 0,
              }}
            >
              <h2 className="font-bold text-3xl tracking-tight md:text-4xl lg:text-5xl">
                {content.cta.title}
              </h2>
              <p className="mx-auto max-w-[700px] text-primary-foreground/80 md:text-xl">
                {content.cta.subtitle}
              </p>
              <div className="mt-4 flex flex-col gap-4 sm:flex-row">
                <Button
                  className="h-12 rounded-full px-8 text-base"
                  size="lg"
                  variant="secondary"
                >
                  {content.cta.buttons.primary.label}
                  <ArrowRightIcon className="ml-2 size-4" />
                </Button>
                <Button
                  className="h-12 rounded-full border-white bg-transparent px-8 text-base text-white hover:bg-white/10"
                  size="lg"
                  variant="outline"
                >
                  {content.cta.buttons.secondary.label}
                </Button>
              </div>
              <p className="mt-4 text-primary-foreground/80 text-sm">
                {content.cta.disclaimer}
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
                  {content.footer.company.logoText}
                </div>
                <span>{content.footer.company.name}</span>
              </div>
              <p className="text-muted-foreground text-sm">
                {content.footer.company.tagline}
              </p>
              <div className="flex gap-4">
                {content.footer.social.map((social) => (
                  <Link
                    className="text-muted-foreground transition-colors hover:text-foreground"
                    key={social.name}
                    to={social.to}
                  >
                    {social.icon === "facebook" && (
                      <svg
                        aria-label="Facebook"
                        className="size-5"
                        fill="none"
                        height="24"
                        role="img"
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
                    )}
                    {social.icon === "twitter" && (
                      <svg
                        aria-label="Twitter"
                        className="size-5"
                        fill="none"
                        height="24"
                        role="img"
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
                    )}
                    {social.icon === "linkedin" && (
                      <svg
                        aria-label="LinkedIn"
                        className="size-5"
                        fill="none"
                        height="24"
                        role="img"
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
                    )}
                    <span className="sr-only">{social.name}</span>
                  </Link>
                ))}
              </div>
            </div>
            <div className="space-y-4">
              <h4 className="font-bold text-sm">Product</h4>
              <ul className="space-y-2 text-sm">
                {content.footer.links.product.map((link) => (
                  <li key={link.label}>
                    <Link
                      className="text-muted-foreground transition-colors hover:text-foreground"
                      {...(link.scrollTo
                        ? {
                            search: {
                              scrollTo: link.scrollTo,
                            },
                          }
                        : {})}
                      to={link.to}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            <div className="space-y-4">
              <h4 className="font-bold text-sm">Resources</h4>
              <ul className="space-y-2 text-sm">
                {content.footer.links.resources.map((link) => (
                  <li key={link.label}>
                    <Link
                      className="text-muted-foreground transition-colors hover:text-foreground"
                      to={link.to}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            <div className="space-y-4">
              <h4 className="font-bold text-sm">Company</h4>
              <ul className="space-y-2 text-sm">
                {content.footer.links.company.map((link) => (
                  <li key={link.label}>
                    <Link
                      className="text-muted-foreground transition-colors hover:text-foreground"
                      to={link.to}
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
          <div className="flex flex-col items-center justify-between gap-4 border-border/40 border-t pt-8 sm:flex-row">
            <p className="text-muted-foreground text-xs">
              {content.footer.legal.copyright}
            </p>
            <div className="flex gap-4">
              {content.footer.legal.links.map((link) => (
                <Link
                  className="text-muted-foreground text-xs transition-colors hover:text-foreground"
                  key={link.label}
                  to={link.to}
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
