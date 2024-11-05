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
    "Beautifully designed components that you can copy and paste into your apps. Accessible. Customizable. Open Source.",
  name: "Arches AI",
  ogImage: "https://ui.shadcn.com/og.jpg",
  routes: [
    {
      href: "/playground",
      Icon: SquareTerminal,
      title: "Playground",
    },
    {
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
          href: "/tools/view",
          Icon: SquareTerminal,
          title: "Explore Tools",
        },
        {
          href: "/tools/runs",
          Icon: PackageCheck,
          title: "Tool Runs",
        },
      ],
      href: "/tools",
      Icon: PackageCheck,
      title: "Tools",
    },
    {
      children: [
        {
          href: "/pipelines/view",
          Icon: Workflow,
          title: "View Pipelines",
        },
      ],
      href: "/pipelines",
      Icon: Workflow,
      title: "Pipelines",
    },
    {
      children: [
        {
          href: "/chatbots/view",
          Icon: Bot,
          title: "Chatbots",
        },
        {
          href: "/chatbots/threads",
          Icon: ListMinus,
          title: "Threads",
        },
      ],
      href: "/chatbots",
      Icon: MessageSquare,
      title: "Chat",
    },
    {
      children: [
        {
          href: "/profile/general",
          Icon: User,
          showInTabs: true,
          title: "Profile",
        },
        {
          href: "/profile/security",
          Icon: KeySquare,
          showInTabs: true,
          title: "Security",
        },
      ],
      href: "/profile",
      Icon: User,
      title: "Account",
    },
    {
      children: [
        {
          href: "/organization/general",
          Icon: Building2,
          showInTabs: true,
          title: "General",
        },
        {
          href: "/organization/billing",
          Icon: CreditCard,
          showInTabs: true,
          title: "Billing",
        },
        {
          href: "/organization/members",
          Icon: Users,
          showInTabs: true,

          title: "Members",
        },
        {
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
