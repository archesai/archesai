import type { SiteRoute } from '@archesai/ui/lib/site-config.interface'

import {
  BookOpen,
  Building2,
  Cpu,
  CpuIcon,
  FileText,
  Image,
  PackageCheck,
  Server,
  Settings2,
  SquareTerminal,
  Tags,
  User,
  Users,
  Volume2,
  Workflow
} from '@archesai/ui/components/custom/icons'

export const siteRoutes: SiteRoute[] = [
  {
    description: 'Try out your tools here.',
    href: '/',
    Icon: SquareTerminal,
    section: 'Home',
    title: 'Dashboard'
  },
  {
    description: 'Browse and manage your artifacts here.',
    href: '/artifacts',
    Icon: Server,
    section: 'Data',
    title: 'Artifacts'
  },
  {
    description: 'Explore and run tools.',
    href: '/tools',
    Icon: CpuIcon,
    section: 'Build',
    title: 'Tools'
  },
  {
    description: 'View and manage your pipelines.',
    href: '/pipelines',
    Icon: Workflow,
    section: 'Build',
    showInTabs: true,
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
        href: '/organization',
        Icon: Building2,
        section: 'Settings',
        showInTabs: true,
        title: 'General'
      },
      {
        description: "View and manage your organization's members.",
        href: '/organization/members',
        Icon: Users,
        section: 'Settings',
        showInTabs: true,
        title: 'Members'
      }
    ],
    href: '/organization',
    Icon: Settings2,
    section: 'Settings',
    title: 'Settings'
  }
]

export const toolBaseIcons: Record<
  | 'create-embeddings'
  | 'extract-text'
  | 'summarize'
  | 'text-to-image'
  | 'text-to-speech',
  SiteRoute['Icon']
> = {
  'create-embeddings': Cpu,
  'extract-text': FileText,
  summarize: BookOpen,
  'text-to-image': Image,
  'text-to-speech': Volume2
}

export const siteMetadata = {
  description:
    'Arches AI is a platform that provides tools to transform data into various forms of content.',
  name: 'Arches AI',
  ogImage: 'https://ui.shadcn.com/og.jpg',
  url: 'https://archesai.com'
}
