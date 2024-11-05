import {
  Bot,
  // Award,
  // Bolt,
  Building2,
  // Compass,
  CreditCard,
  // DollarSign,
  // HelpCircle,
  // HomeIcon,
  KeySquare,
  ListMinus,
  // Layers,
  Lock,
  MessageSquare,
  PackageCheck,
  Server,
  Settings2,
  SquareTerminal,
  User,
  Users,
  Workflow,
} from "lucide-react";

export const siteConfig = {
  description:
    "Arches AI is a platform that provides tools to transform data into various forms of content.",
  name: "Arches AI",
  ogImage: "https://ui.shadcn.com/og.jpg",
  routes: [
    {
      description: "Try out your tools here. If you like them, save them.",
      href: "/playground",
      Icon: SquareTerminal,
      title: "Playground",
    },
    {
      description: "Browse and manage your content here.",
      // children: [
      //   {
      //     href: "/content/view",
      //     Icon: Server,
      //     title: "Browse",
      //   },
      // ],
      href: "/content/view",
      Icon: Server,
      title: "Content",
    },
    {
      children: [
        {
          description: "View and manage your tools.",
          href: "/tools/view",
          Icon: SquareTerminal,
          title: "Explore Tools",
        },
        {
          description:
            "View your tool runs. You can track progress, retry or cancel runs.",
          href: "/tools/runs",
          Icon: PackageCheck,
          title: "Tool Runs",
        },
      ],
      description: "Explore and run tools.",
      href: "/tools",
      Icon: PackageCheck,
      title: "Tools",
    },
    {
      children: [
        {
          description: "View and manage your pipelines.",
          href: "/pipelines/view",
          Icon: Workflow,
          title: "View Pipelines",
        },
      ],
      description: "Create and manage pipelines.",
      href: "/pipelines",
      Icon: Workflow,
      title: "Pipelines",
    },
    {
      children: [
        {
          href: "/chatbots/chat",
          Icon: Bot,
          title: "New Thread",
        },
        {
          description: "View and manage chatbot threads.",
          href: "/chatbots/threads",
          Icon: ListMinus,
          title: "History",
        },
      ],
      description: "Create and manage chatbots.",
      href: "/chatbots",
      Icon: MessageSquare,
      title: "Chat",
    },
    {
      children: [
        {
          description: "Update your profile information.",
          href: "/profile/general",
          Icon: User,
          showInTabs: true,
          title: "Profile",
        },
        {
          description: "Update your security settings.",
          href: "/profile/security",
          Icon: KeySquare,
          showInTabs: true,
          title: "Security",
        },
      ],
      description: "View your profile information.",
      href: "/profile",
      Icon: User,
      title: "Account",
    },
    {
      children: [
        {
          description:
            "View and update your organization's general information.",
          href: "/organization/general",
          Icon: Building2,
          showInTabs: true,
          title: "General",
        },
        {
          description:
            "View and update your organization's billing information. Upgrade your plan.",
          href: "/organization/billing",
          Icon: CreditCard,
          showInTabs: true,
          title: "Billing",
        },
        {
          description: "View and manage your organization's members.",
          href: "/organization/members",
          Icon: Users,
          showInTabs: true,
          title: "Members",
        },
        {
          description: "View and manage your organization's API tokens.",
          href: "/organization/api-tokens",
          Icon: Lock,
          showInTabs: true,
          title: "API Tokens",
        },
      ],
      href: "/organization",
      Icon: Settings2,
      title: "Settings",
    },
  ],
  url: "https://archesai.com",
};

export type SiteConfig = typeof siteConfig;
