import {
  // Bot,
  // Award,
  // Bolt,
  Building2,
  CpuIcon,
  // Compass,
  CreditCard,
  // DollarSign,
  // HelpCircle,
  // HomeIcon,
  KeySquare,
  // ListMinus,
  // Layers,
  Lock,
  MessageSquareIcon,
  PackageCheck,
  Server,
  Settings2,
  SquareTerminal,
  Tags,
  User,
  Users,
  Workflow
} from 'lucide-react'
import { BookOpen, Cpu, FileText, Image, Volume2 } from 'lucide-react'

export const siteConfig = {
  description:
    'Arches AI is a platform that provides tools to transform data into various forms of content.',
  name: 'Arches AI',
  ogImage: 'https://ui.shadcn.com/og.jpg',
  routes: [
    {
      description: 'Try out your tools here.',
      href: '/playground',
      Icon: SquareTerminal,
      section: 'Build',
      title: 'Playground'
    },
    {
      description: 'Browse and manage your content here.',
      href: '/content/view',
      Icon: Server,
      section: 'Data',
      title: 'Content'
    },
    {
      description: 'Explore and run tools.',
      href: '/tools',
      Icon: CpuIcon,
      section: 'Build',
      title: 'Tools'
    },
    {
      children: [
        {
          description: 'View and manage your pipelines.',
          href: '/pipelines',
          Icon: Workflow,
          showInTabs: true,
          title: 'Pipelines'
        },
        {
          description: 'Create a new pipeline.',
          href: '/pipelines/create',
          Icon: Workflow,
          showInTabs: false,
          title: 'Create'
        }
      ],
      description: 'Create and manage pipelines.',
      href: '/pipelines',
      Icon: Workflow,
      section: 'Build',
      title: 'Pipelines'
    },
    {
      description: 'View your previous runs.',
      href: '/runs',
      Icon: PackageCheck,
      section: 'Build',
      title: 'History'
    },
    {
      description: 'Create and manage labels.',
      href: '/labels',
      Icon: Tags,
      section: 'Data',
      title: 'Labels'
    },
    {
      description: 'Chat with your data and tools.',
      href: '/chat',
      Icon: MessageSquareIcon,
      section: 'Chat',
      title: 'Chat'
    },
    {
      children: [
        {
          description: 'Update your profile information.',
          href: '/profile/general',
          Icon: User,
          showInTabs: true,
          title: 'Profile'
        },
        {
          description: 'Update your security settings.',
          href: '/profile/security',
          Icon: KeySquare,
          showInTabs: true,
          title: 'Security'
        }
      ],
      description: 'View your profile information.',
      href: '/profile/general',
      Icon: User,
      section: 'Settings',
      title: 'Account'
    },
    {
      children: [
        {
          description:
            "View and update your organization's general information.",
          href: '/organization/general',
          Icon: Building2,
          showInTabs: true,
          title: 'General'
        },
        {
          description:
            "View and update your organization's billing information. Upgrade your plan.",
          href: '/organization/billing',
          Icon: CreditCard,
          showInTabs: true,
          title: 'Billing'
        },
        {
          description: "View and manage your organization's members.",
          href: '/organization/members',
          Icon: Users,
          showInTabs: true,
          title: 'Members'
        },
        {
          description: "View and manage your organization's API tokens.",
          href: '/organization/api-tokens',
          Icon: Lock,
          showInTabs: true,
          title: 'API Tokens'
        }
      ],
      href: '/organization/general',
      Icon: Settings2,
      section: 'Settings',
      title: 'Settings'
    }
  ],
  toolBaseIcons: {
    'create-embeddings': Cpu,
    'extract-text': FileText,
    summarize: BookOpen,
    'text-to-image': Image,
    'text-to-speech': Volume2
  } as Record<string, any>,
  url: 'https://archesai.com'
}

export const getMetadata = (href: string) => {
  // do this without using flatMap
  const childRoutes = siteConfig.routes
    .map((route) => route.children)
    .filter((children) => children)

  const allRoutes = siteConfig.routes.concat(
    childRoutes.reduce((acc, val) => acc.concat(val), [] as any)
  )
  return {
    description: allRoutes.find((route) => route.href === href)?.description,
    title:
      allRoutes.find((route) => route.href === href)?.title + ' | Arches AI'
  }
}

export type SiteConfig = typeof siteConfig
