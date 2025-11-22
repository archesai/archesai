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
        label: "Start Saving Now",
      },
      secondary: {
        label: "See Live Demo",
      },
    },
    disclaimer:
      "30-day free trial. No credit card. Setup in 5 minutes. Cancel anytime.",
    subtitle:
      "Every day without Arches costs you thousands in inefficiency. Join 10,000+ companies already saving millions.",
    title: "Stop Losing Money on Manual Processes",
  },
  faq: {
    badge: "FAQ",
    questions: [
      {
        answer:
          "Most customers see immediate results. On average, companies report 50% time savings on day one, and 10x ROI within the first week. Our AI begins optimizing your workflows from the moment you connect your data, with full optimization typically achieved within 48 hours.",
        question: "How quickly can I see results with Arches?",
      },
      {
        answer:
          "Absolutely not. Arches is designed for business users, not engineers. If you can use email, you can use Arches. Our AI handles all the technical complexity behind the scenes. We also provide free onboarding and training for all new customers.",
        question: "Do I need technical expertise to use Arches?",
      },
      {
        answer:
          "Building similar capabilities in-house would require a team of 10-15 engineers and cost $2-3M annually. With Arches, you get enterprise-grade AI infrastructure for less than the cost of a single junior developer. Plus, you're live in 5 minutes instead of 18 months.",
        question: "How does Arches compare to building in-house?",
      },
      {
        answer:
          "Your data never leaves our SOC 2 Type II certified infrastructure. We use military-grade AES-256 encryption, and our security is audited quarterly by independent firms. We're GDPR, CCPA, and HIPAA compliant. Your data is safer with us than on your own servers.",
        question: "What happens to my data? Is it secure?",
      },
      {
        answer:
          "We process billions of data points daily for Fortune 500 companies. Our infrastructure auto-scales to handle any volume - from 10 to 10 billion operations. Companies like Microsoft and Amazon trust us with their mission-critical workflows.",
        question: "Can Arches handle my scale?",
      },
      {
        answer:
          "You can cancel anytime with one click. No questions asked, no hidden fees. We'll even help you export all your data and provide 30 days of free access to ensure a smooth transition. But with our 98% retention rate, we're confident you'll want to stay.",
        question: "What if I need to cancel?",
      },
      {
        answer:
          "Our proprietary AI models are trained on petabytes of business data to understand patterns and optimize workflows. The AI learns from your specific use cases and improves continuously. Within a week, it understands your business better than a consultant who's been studying it for months.",
        question: "How does the AI actually work?",
      },
      {
        answer:
          "Our average customer sees 10-15x ROI within 90 days. This comes from: 70% reduction in processing time, 90% fewer errors, 50% less staff time on repetitive tasks, and the ability to scale without hiring. We guarantee at least 5x ROI or your money back.",
        question: "What ROI can I expect?",
      },
    ],
    subtitle:
      "Got questions? We've got answers. Can't find what you're looking for? Our team responds in under 2 minutes.",
    title: "Everything You Need to Know",
  },
  features: {
    badge: "Why Arches",
    features: [
      {
        description:
          "Our proprietary AI models learn from your patterns and automate complex workflows with 99.9% accuracy. Save 30+ hours per week on repetitive tasks.",
        icon: "ZapIcon",
        title: "AI-Powered Automation",
      },
      {
        description:
          "Get instant insights with predictive analytics that spots trends before they happen. Make data-driven decisions 5x faster than traditional methods.",
        icon: "BarChartIcon",
        title: "Real-Time Intelligence",
      },
      {
        description:
          "Bring your entire team together in one intelligent platform. Real-time collaboration with AI-assisted communication and smart task distribution.",
        icon: "UsersIcon",
        title: "Unified Workspace",
      },
      {
        description:
          "SOC 2 Type II certified with military-grade encryption. Your data is protected by the same security standards used by Fortune 500 companies.",
        icon: "ShieldIcon",
        title: "Bank-Level Security",
      },
      {
        description:
          "Connect instantly with Salesforce, Slack, Microsoft 365, and 1000+ other tools. Our universal API adapts to any workflow in minutes.",
        icon: "LayersIcon",
        title: "1000+ Integrations",
      },
      {
        description:
          "Dedicated success manager, 24/7 priority support, and 99.99% uptime SLA. Average response time under 2 minutes.",
        icon: "StarIcon",
        title: "White-Glove Support",
      },
    ],
    subtitle:
      "Built from the ground up to handle complex workflows at scale. Our AI-native architecture processes millions of data points in seconds, not hours.",
    title: "The Future of Intelligent Data Processing",
  },
  footer: {
    company: {
      logoText: "A",
      name: "Arches",
      tagline:
        "The world's most advanced AI-powered data processing platform. Trusted by industry leaders to handle mission-critical workflows at scale.",
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
    badge: "AI-Powered Platform",
    benefits: [
      "No credit card required",
      "30-day free trial",
      "Setup in 5 minutes",
    ],
    buttons: {
      primary: {
        label: "Start Free Trial",
      },
      secondary: {
        label: "Watch 2-Min Demo",
      },
    },
    image: {
      alt: "Arches intelligent dashboard showing real-time data processing",
      src: "https://cdn.dribbble.com/userupload/12302729/file/original-fa372845e394ee85bebe0389b9d86871.png?resize=1504x1128&vertical=center",
    },
    subtitle:
      "Arches combines cutting-edge AI with powerful automation to process data 10x faster. Join industry leaders who've reduced operational costs by 40% while scaling effortlessly.",
    title: "Transform Your Data Into Intelligent Workflows",
  },
  howItWorks: {
    badge: "Get Started",
    steps: [
      {
        description:
          "One-click integration with your existing tools. Our AI automatically maps and organizes your data structure.",
        step: "01",
        title: "Connect Your Data",
      },
      {
        description:
          "Our AI analyzes your processes and creates custom automation rules. See optimization suggestions in real-time.",
        step: "02",
        title: "AI Learns Your Workflow",
      },
      {
        description:
          "Watch as tasks that took hours complete in seconds. Scale from 10 to 10,000 operations without changing a thing.",
        step: "03",
        title: "Scale Effortlessly",
      },
    ],
    subtitle:
      "Join thousands of companies that went live in under 5 minutes. No technical expertise required.",
    title: "From Zero to AI-Powered in Minutes",
  },
  logos: {
    logos: [
      {
        alt: "Microsoft",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Google",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Amazon",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Meta",
        src: "/placeholder-logo.svg",
      },
      {
        alt: "Tesla",
        src: "/placeholder-logo.svg",
      },
    ],
    title: "Powering the world's most innovative companies",
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
      "Start with 30 days free. Most customers see 10x ROI within the first week.",
    tabs: {
      annually: {
        label: "Annually",
        plans: [
          {
            cta: "Start Free Trial",
            description: "For ambitious teams ready to scale.",
            features: [
              "Process up to 1M data points/month",
              "10 AI automation workflows",
              "100GB high-speed storage",
              "5 team members included",
              "Standard integrations (50+)",
              "Email & chat support",
              "99.9% uptime SLA",
            ],
            name: "Growth",
            price: "$224",
          },
          {
            cta: "Start Free Trial",
            description: "For companies experiencing rapid growth.",
            features: [
              "Process up to 10M data points/month",
              "Unlimited AI workflows",
              "1TB high-speed storage",
              "25 team members included",
              "Premium integrations (500+)",
              "Priority 24/7 support",
              "Custom AI model training",
              "99.99% uptime SLA",
              "Advanced security features",
            ],
            name: "Scale",
            popular: true,
            price: "$749",
          },
          {
            cta: "Contact Sales",
            description: "For industry leaders with mission-critical needs.",
            features: [
              "Unlimited data processing",
              "Unlimited everything",
              "Dedicated infrastructure",
              "Unlimited team members",
              "Custom integrations",
              "Dedicated success team",
              "On-premise deployment option",
              "100% uptime SLA",
              "Custom AI models",
              "White-label options",
              "Priority roadmap input",
            ],
            name: "Enterprise",
            price: "Custom",
          },
        ],
        savingsText: "Save 25%",
      },
      monthly: {
        label: "Monthly",
        plans: [
          {
            cta: "Start Free Trial",
            description: "For ambitious teams ready to scale.",
            features: [
              "Process up to 1M data points/month",
              "10 AI automation workflows",
              "100GB high-speed storage",
              "5 team members included",
              "Standard integrations (50+)",
              "Email & chat support",
              "99.9% uptime SLA",
            ],
            name: "Growth",
            price: "$299",
          },
          {
            cta: "Start Free Trial",
            description: "For companies experiencing rapid growth.",
            features: [
              "Process up to 10M data points/month",
              "Unlimited AI workflows",
              "1TB high-speed storage",
              "25 team members included",
              "Premium integrations (500+)",
              "Priority 24/7 support",
              "Custom AI model training",
              "99.99% uptime SLA",
              "Advanced security features",
            ],
            name: "Scale",
            popular: true,
            price: "$999",
          },
          {
            cta: "Contact Sales",
            description: "For industry leaders with mission-critical needs.",
            features: [
              "Unlimited data processing",
              "Unlimited everything",
              "Dedicated infrastructure",
              "Unlimited team members",
              "Custom integrations",
              "Dedicated success team",
              "On-premise deployment option",
              "100% uptime SLA",
              "Custom AI models",
              "White-label options",
              "Priority roadmap input",
            ],
            name: "Enterprise",
            price: "Custom",
          },
        ],
      },
    },
    title: "ROI-Positive From Day One",
  },
  testimonials: {
    badge: "Success Stories",
    subtitle:
      "Join 10,000+ companies experiencing unprecedented growth with Arches.",
    testimonials: [
      {
        author: "Sarah Chen",
        quote:
          "Arches reduced our data processing time by 93%. What used to take our team 8 hours now completes in 30 minutes. The ROI was immediate - we saved $2M in the first quarter alone.",
        rating: 5,
        role: "CTO, Fortune 500 Tech Company",
      },
      {
        author: "Marcus Williams",
        quote:
          "The AI predictions are scary accurate. We prevented 3 major supply chain disruptions last month alone. Arches paid for itself 10x over in the first week.",
        rating: 5,
        role: "VP Operations, Global E-commerce Leader",
      },
      {
        author: "Dr. Emily Thompson",
        quote:
          "Processing genomic data that took weeks now takes hours. Arches accelerated our research timeline by 18 months. This is the future of scientific computing.",
        rating: 5,
        role: "Head of Research, BioTech Unicorn",
      },
      {
        author: "David Park",
        quote:
          "We scaled from 100 to 10,000 customers without hiring a single data analyst. Arches handles everything. It's like having a team of 50 engineers for the price of 1.",
        rating: 5,
        role: "CEO, Hypergrowth SaaS Startup",
      },
      {
        author: "Lisa Anderson",
        quote:
          "Compliance reporting that took 2 weeks now generates in real-time. We've eliminated 100% of manual errors and saved $5M annually. Best investment we've ever made.",
        rating: 5,
        role: "CFO, Financial Services Giant",
      },
      {
        author: "Roberto Silva",
        quote:
          "Arches predicted equipment failures 14 days in advance with 98% accuracy. We've reduced downtime by 87% and increased production efficiency by 45%.",
        rating: 5,
        role: "Director of Innovation, Manufacturing Leader",
      },
    ],
    title: "Results That Speak for Themselves",
  },
};
