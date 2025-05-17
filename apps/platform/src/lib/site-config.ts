import {
  BookOpen,
  Building2,
  Cpu,
  CpuIcon,
  CreditCard,
  FileText,
  Image,
  KeySquare,
  Lock,
  MessageSquareIcon,
  PackageCheck,
  Server,
  Settings2,
  SquareTerminal,
  Tags,
  User,
  Users,
  Volume2,
  Workflow
} from 'lucide-react'

import type {
  RouteIcon,
  SiteRoute
} from '@archesai/ui/lib/site-config.interface'

export const siteRoutes: SiteRoute[] = [
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
        section: 'Build',
        showInTabs: true,
        title: 'Pipelines'
      },
      {
        description: 'Create a new pipeline.',
        href: '/pipelines/create',
        Icon: Workflow,
        section: 'Build',
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
        section: 'Settings',
        showInTabs: true,
        title: 'Profile'
      },
      {
        description: 'Update your security settings.',
        href: '/profile/security',
        Icon: KeySquare,
        section: 'Settings',
        showInTabs: true,
        title: 'Security'
      }
    ],
    description: 'View your profile information.',
    href: '/profile',
    Icon: User,
    section: 'Settings',
    title: 'Account'
  },
  {
    children: [
      {
        description: "View and update your organization's general information.",
        href: '/organization/general',
        Icon: Building2,
        section: 'Settings',
        showInTabs: true,
        title: 'General'
      },
      {
        description:
          "View and update your organization's billing information. Upgrade your plan.",
        href: '/organization/billing',
        Icon: CreditCard,
        section: 'Settings',
        showInTabs: true,
        title: 'Billing'
      },
      {
        description: "View and manage your organization's members.",
        href: '/organization/members',
        Icon: Users,
        section: 'Settings',
        showInTabs: true,
        title: 'Members'
      },
      {
        description: "View and manage your organization's API tokens.",
        href: '/organization/api-tokens',
        Icon: Lock,
        section: 'Settings',
        showInTabs: true,
        title: 'API Tokens'
      }
    ],
    href: '/organization',
    Icon: Settings2,
    section: 'Settings',
    title: 'Settings'
  }
]

export const toolBaseIcons = {
  'create-embeddings': Cpu,
  'extract-text': FileText,
  summarize: BookOpen,
  'text-to-image': Image,
  'text-to-speech': Volume2
} satisfies Record<string, RouteIcon>

export const siteMetadata = {
  description:
    'Arches AI is a platform that provides tools to transform data into various forms of content.',
  name: 'Arches AI',
  ogImage: 'https://ui.shadcn.com/og.jpg',
  url: 'https://archesai.com'
}
