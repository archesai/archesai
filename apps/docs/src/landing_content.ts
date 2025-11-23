export interface LandingContent {
  navigation: NavigationContent;
  hero: HeroContent;
  logos: LogosContent;
  features: FeaturesContent;
  howItWorks: HowItWorksContent;
  testimonials: TestimonialsContent;
  pricing: PricingContent;
  faq: FAQContent;
  cta: CTAContent;
  footer: FooterContent;
}

interface NavigationContent {
  links: Array<{
    label: string;
    scrollTo?: string;
    to: string;
  }>;
  buttons: {
    login: {
      label: string;
      to: string;
    };
    getStarted: {
      label: string;
    };
  };
}

interface HeroContent {
  badge: string;
  title: string;
  subtitle: string;
  buttons: {
    primary: {
      label: string;
    };
    secondary: {
      label: string;
    };
  };
  benefits: string[];
  image: {
    src: string;
    alt: string;
  };
}

interface LogosContent {
  title: string;
  logos: Array<{
    src: string;
    alt: string;
  }>;
}

interface FeaturesContent {
  badge: string;
  title: string;
  subtitle: string;
  features: Array<{
    title: string;
    description: string;
    icon: string;
  }>;
}

interface HowItWorksContent {
  badge: string;
  title: string;
  subtitle: string;
  steps: Array<{
    step: string;
    title: string;
    description: string;
  }>;
}

interface TestimonialsContent {
  badge: string;
  title: string;
  subtitle: string;
  testimonials: Array<{
    author: string;
    role: string;
    quote: string;
    rating: number;
  }>;
}

interface PricingPlan {
  name: string;
  price: string;
  description: string;
  features: string[];
  cta: string;
  popular?: boolean;
}

interface PricingContent {
  badge: string;
  title: string;
  subtitle: string;
  tabs: {
    monthly: {
      label: string;
      plans: PricingPlan[];
    };
    annually: {
      label: string;
      savingsText: string;
      plans: PricingPlan[];
    };
  };
}

interface FAQContent {
  badge: string;
  title: string;
  subtitle: string;
  questions: Array<{
    question: string;
    answer: string;
  }>;
}

interface CTAContent {
  title: string;
  subtitle: string;
  buttons: {
    primary: {
      label: string;
    };
    secondary: {
      label: string;
    };
  };
  disclaimer: string;
}

interface FooterContent {
  company: {
    name: string;
    tagline: string;
    logoText: string;
  };
  social: Array<{
    name: string;
    to: string;
    icon: string;
  }>;
  links: {
    product: Array<{
      label: string;
      to: string;
      scrollTo?: string;
    }>;
    resources: Array<{
      label: string;
      to: string;
    }>;
    company: Array<{
      label: string;
      to: string;
    }>;
  };
  legal: {
    copyright: string;
    links: Array<{
      label: string;
      to: string;
    }>;
  };
}

export const defaultContent: LandingContent = {
  cta: {
    buttons: {
      primary: {
        label: "Start Building Now",
      },
      secondary: {
        label: "View Documentation",
      },
    },
    disclaimer:
      "Free and open source. No credit card required. Generate your first app in 5 minutes.",
    subtitle:
      "Join thousands of developers who are shipping production apps 10x faster with OpenAPI-driven development.",
    title: "Ready to Build Your Next App?",
  },
  faq: {
    badge: "FAQ",
    questions: [
      {
        answer:
          "In about 5 minutes! Once you install Arches and create your OpenAPI specification, you can generate a complete application instantly. The generated code includes everything you need: models, controllers, database migrations, authentication, and even a TypeScript client SDK.",
        question: "How quickly can I generate my first app?",
      },
      {
        answer:
          "Basic programming knowledge is helpful but not required for simple apps. Arches generates production-ready code that follows best practices. You'll need to understand your API design (OpenAPI) and potentially add custom business logic, but the heavy lifting is done for you.",
        question: "Do I need to be an expert developer to use Arches?",
      },
      {
        answer:
          "Building a production-ready app manually can take weeks or months. With Arches, you get the same result in minutes. We estimate Arches saves 80-90% of boilerplate coding time, letting you focus on your unique business logic instead of repetitive infrastructure code.",
        question: "How does Arches compare to coding from scratch?",
      },
      {
        answer:
          "Arches is completely open source under the AGPLv3 license. Your generated code belongs to you. We don't store any of your specifications or generated code - everything runs locally on your machine. You have full control over your applications.",
        question: "What happens to my code? Who owns it?",
      },
      {
        answer:
          "Absolutely! Arches generates clean, well-structured code that follows industry best practices. The generated code is designed to scale from prototype to production. Many teams use Arches for hackathons, MVPs, and production applications serving millions of users.",
        question: "Can I use Arches for production applications?",
      },
      {
        answer:
          "Arches currently generates Go backends with full TypeScript/JavaScript client SDKs. Python and Node.js backend support is coming soon. The generated code includes Docker configurations and Kubernetes manifests for easy deployment anywhere.",
        question: "What languages and frameworks does Arches support?",
      },
      {
        answer:
          "Yes! The generated code is meant to be extended. Arches uses special markers (x-codegen annotations) to let you customize generation behavior. You can add custom business logic, override generated methods, and even provide your own templates for complete control.",
        question: "Can I customize the generated code?",
      },
      {
        answer:
          "Arches generates comprehensive auth systems including JWT tokens, email/password login, OAuth integration, magic links, session management, and role-based access control. Just define your security requirements in your OpenAPI spec and Arches handles the implementation.",
        question: "How does authentication work in generated apps?",
      },
    ],
    subtitle:
      "Got questions? We've got answers. Check our documentation or open an issue on GitHub for more help.",
    title: "Frequently Asked Questions",
  },
  features: {
    badge: "Why Arches",
    features: [
      {
        description:
          "Generate complete applications from OpenAPI specifications in seconds. Models, controllers, database layers, authentication - everything you need, production-ready.",
        icon: "ZapIcon",
        title: "OpenAPI to Full App",
      },
      {
        description:
          "Clean, maintainable Go code following best practices. Includes comprehensive error handling, validation, and testing infrastructure out of the box.",
        icon: "BarChartIcon",
        title: "Production-Ready Code",
      },
      {
        description:
          "Generated TypeScript/JavaScript client SDKs with full type safety. Automatic validation, error handling, and API documentation included.",
        icon: "UsersIcon",
        title: "Type-Safe Client SDKs",
      },
      {
        description:
          "JWT authentication, OAuth providers, magic links, session management, and RBAC. Enterprise-grade security patterns built into every generated app.",
        icon: "ShieldIcon",
        title: "Built-in Authentication",
      },
      {
        description:
          "Docker, Kubernetes, and Helm charts generated automatically. Deploy anywhere with production-ready configurations and health checks.",
        icon: "LayersIcon",
        title: "Cloud-Native Deployment",
      },
      {
        description:
          "Hot reload during development, extensible templates, custom annotations, and powerful CLI tools. Built by developers, for developers.",
        icon: "StarIcon",
        title: "Developer Experience",
      },
    ],
    subtitle:
      "Transform your OpenAPI specifications into complete, production-ready applications. Skip the boilerplate and focus on your business logic.",
    title: "Code Generation That Actually Works",
  },
  footer: {
    company: {
      logoText: "A",
      name: "Arches",
      tagline:
        "Open-source code generation platform that transforms OpenAPI specifications into production-ready applications. Build faster, ship sooner.",
    },
    legal: {
      copyright: `© ${new Date().getFullYear()} Arches. All rights reserved.`,
      links: [
        {
          label: "Privacy Policy",
          to: "/",
        },
        {
          label: "Terms of Service",
          to: "/",
        },
        {
          label: "Cookie Policy",
          to: "/",
        },
        {
          label: "GDPR",
          to: "/",
        },
        {
          label: "Security",
          to: "/",
        },
      ],
    },
    links: {
      company: [
        {
          label: "About Us",
          to: "/",
        },
        {
          label: "Careers",
          to: "/",
        },
        {
          label: "Press Kit",
          to: "/",
        },
        {
          label: "Investors",
          to: "/",
        },
        {
          label: "Contact",
          to: "/",
        },
        {
          label: "Partners",
          to: "/",
        },
      ],
      product: [
        {
          label: "Features",
          scrollTo: "features",
          to: "/",
        },
        {
          label: "Pricing",
          scrollTo: "pricing",
          to: "/",
        },
        {
          label: "Integrations",
          to: "/",
        },
        {
          label: "API Documentation",
          to: "/",
        },
        {
          label: "Security",
          to: "/",
        },
        {
          label: "Roadmap",
          to: "/",
        },
      ],
      resources: [
        {
          label: "Documentation",
          to: "/",
        },
        {
          label: "Video Tutorials",
          to: "/",
        },
        {
          label: "Case Studies",
          to: "/",
        },
        {
          label: "Blog",
          to: "/",
        },
        {
          label: "Community",
          to: "/",
        },
        {
          label: "System Status",
          to: "/",
        },
      ],
    },
    social: [
      {
        icon: "linkedin",
        name: "LinkedIn",
        to: "/",
      },
      {
        icon: "twitter",
        name: "Twitter",
        to: "/",
      },
      {
        icon: "facebook",
        name: "Facebook",
        to: "/",
      },
    ],
  },
  hero: {
    badge: "OpenAPI-Driven Development",
    benefits: ["Open source", "Generate in seconds", "Production ready"],
    buttons: {
      primary: {
        label: "Get Started",
      },
      secondary: {
        label: "View Docs",
      },
    },
    image: {
      alt: "Arches code generation showing OpenAPI to application transformation",
      src: "https://cdn.dribbble.com/userupload/12302729/file/original-fa372845e394ee85bebe0389b9d86871.png?resize=1504x1128&vertical=center",
    },
    subtitle:
      "Arches transforms your OpenAPI specifications into complete, production-ready applications. Skip weeks of boilerplate coding and focus on what makes your app unique.",
    title: "From API Spec to Production App in Minutes",
  },
  howItWorks: {
    badge: "Get Started",
    steps: [
      {
        description:
          "Design your API using OpenAPI 3.0/3.1 specification. Define endpoints, models, and authentication requirements.",
        step: "01",
        title: "Define Your API",
      },
      {
        description:
          "Run the Arches CLI to transform your specification into a complete application with all the boilerplate handled.",
        step: "02",
        title: "Generate Your App",
      },
      {
        description:
          "Add your business logic to the generated handlers. Deploy with included Docker/Kubernetes configs.",
        step: "03",
        title: "Customize & Deploy",
      },
    ],
    subtitle:
      "Join thousands of developers building production apps in minutes, not months.",
    title: "Three Steps to Production",
  },
  logos: {
    logos: [
      {
        alt: "Go",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "OpenAPI",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "PostgreSQL",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Docker",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Kubernetes",
        src: "/placeholder-logo.svg",
      },
    ],
    title: "Built with modern, battle-tested technologies",
  },
  navigation: {
    buttons: {
      getStarted: {
        label: "Get Started",
      },
      login: {
        label: "Log in",
        to: "/",
      },
    },
    links: [
      {
        label: "Features",
        scrollTo: "features",
        to: "/",
      },
      {
        label: "Testimonials",
        scrollTo: "testimonials",
        to: "/",
      },
      {
        label: "Pricing",
        scrollTo: "pricing",
        to: "/",
      },
      {
        label: "FAQ",
        scrollTo: "faq",
        to: "/",
      },
    ],
  },
  pricing: {
    badge: "Pricing",
    subtitle:
      "Arches is free and open source. Use it locally or deploy anywhere.",
    tabs: {
      annually: {
        label: "Open Source",
        plans: [
          {
            cta: "Get Started",
            description: "Everything you need to build production apps.",
            features: [
              "Unlimited code generation",
              "All languages and frameworks",
              "Complete source code access",
              "Community support",
              "Docker & Kubernetes configs",
              "Authentication systems",
              "Client SDK generation",
            ],
            name: "Community",
            price: "Free",
          },
          {
            cta: "View Docs",
            description: "Commercial support and enterprise features.",
            features: [
              "Everything in Community",
              "Priority support",
              "Custom templates",
              "Training and onboarding",
              "SLA guarantees",
              "Custom integrations",
              "White-label options",
              "Dedicated success manager",
              "Priority feature requests",
            ],
            name: "Enterprise",
            popular: true,
            price: "Contact",
          },
          {
            cta: "Learn More",
            description: "Managed cloud service coming soon.",
            features: [
              "Hosted code generation",
              "Team collaboration",
              "Version control integration",
              "CI/CD pipelines",
              "Private template registry",
              "Advanced analytics",
              "Automated deployments",
              "99.99% uptime SLA",
              "24/7 support",
            ],
            name: "Cloud (Coming Soon)",
            price: "TBD",
          },
        ],
        savingsText: "",
      },
      monthly: {
        label: "Self-Hosted",
        plans: [
          {
            cta: "Get Started",
            description: "Everything you need to build production apps.",
            features: [
              "Unlimited code generation",
              "All languages and frameworks",
              "Complete source code access",
              "Community support",
              "Docker & Kubernetes configs",
              "Authentication systems",
              "Client SDK generation",
            ],
            name: "Community",
            price: "Free",
          },
          {
            cta: "View Docs",
            description: "Commercial support and enterprise features.",
            features: [
              "Everything in Community",
              "Priority support",
              "Custom templates",
              "Training and onboarding",
              "SLA guarantees",
              "Custom integrations",
              "White-label options",
              "Dedicated success manager",
              "Priority feature requests",
            ],
            name: "Enterprise",
            popular: true,
            price: "Contact",
          },
          {
            cta: "Learn More",
            description: "Managed cloud service coming soon.",
            features: [
              "Hosted code generation",
              "Team collaboration",
              "Version control integration",
              "CI/CD pipelines",
              "Private template registry",
              "Advanced analytics",
              "Automated deployments",
              "99.99% uptime SLA",
              "24/7 support",
            ],
            name: "Cloud (Coming Soon)",
            price: "TBD",
          },
        ],
      },
    },
    title: "Free & Open Source",
  },
  testimonials: {
    badge: "Developer Stories",
    subtitle: "See what developers are building with Arches.",
    testimonials: [
      {
        author: "Sarah Chen",
        quote:
          "Arches cut our API development time from 3 weeks to 3 hours. The generated code is clean, well-tested, and production-ready. It's exactly what we would have written ourselves, just 100x faster.",
        rating: 5,
        role: "Senior Engineer, YC Startup",
      },
      {
        author: "Marcus Williams",
        quote:
          "We use Arches for all our microservices now. Define the OpenAPI spec, generate the code, add business logic, deploy. What used to take a sprint now takes a day.",
        rating: 5,
        role: "Platform Lead, Fortune 500",
      },
      {
        author: "Emily Thompson",
        quote:
          "The authentication system Arches generates is more comprehensive than what we built manually. JWT, OAuth, magic links - it's all there and properly secured. Saved us months of work.",
        rating: 5,
        role: "Full Stack Developer",
      },
      {
        author: "David Park",
        quote:
          "Built our entire MVP with Arches in a weekend. The generated TypeScript client SDK was a game-changer for our frontend team. We launched 2 months ahead of schedule.",
        rating: 5,
        role: "CTO, SaaS Startup",
      },
      {
        author: "Lisa Anderson",
        quote:
          "Migrated 12 legacy services to modern Go microservices using Arches. The consistency across all services is incredible. Onboarding new developers is now trivial.",
        rating: 5,
        role: "Engineering Manager",
      },
      {
        author: "Roberto Silva",
        quote:
          "Arches is part of our standard stack now. Every new API starts with an OpenAPI spec and Arches generation. It enforces best practices and eliminates entire categories of bugs.",
        rating: 5,
        role: "Principal Engineer",
      },
    ],
    title: "Loved by Developers Worldwide",
  },
};
