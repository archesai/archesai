import {
  // Award,
  // Bolt,
  Building2,
  // Compass,
  CreditCard,
  // DollarSign,
  // HelpCircle,
  // HomeIcon,
  KeySquare,
  // Layers,
  Lock,
  MessageSquare,
  PackageCheck,
  Server,
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
      children: [
        {
          href: "/chatbots/view",
          Icon: MessageSquare,
          title: "View",
        },
      ],
      href: "/chatbots",
      Icon: MessageSquare,
      title: "Chatbots",
    },
    {
      children: [
        {
          href: "/content/view",
          Icon: Server,
          title: "View",
        },
        {
          href: "/content/single",
          Icon: Server,
          title: "View Single Content",
        },
      ],
      href: "/content",
      Icon: Server,
      title: "Content",
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
      Icon: Building2,
      title: "Organization",
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
          href: "/profile/general",
          Icon: User,
          showInTabs: true,
          title: "General",
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
      title: "Profile",
    },
    {
      children: [
        {
          href: "/tools/view",
          Icon: PackageCheck,
          title: "View Tools",
        },
        {
          href: "/tools/runs",
          Icon: PackageCheck,
          title: "View Tool Runs",
        },
      ],
      href: "/tools",
      Icon: PackageCheck,
      title: "Tools",
    },
  ],
  url: "https://archesai.com",
};

export type SiteConfig = typeof siteConfig;
