// config/site.ts
import {
  Award,
  Bolt,
  Building2,
  CloudUploadIcon,
  Compass,
  DollarSign,
  HelpCircle,
  HomeIcon,
  Image,
  Key,
  KeySquare,
  Layers,
  MessageSquare,
  Server,
  User,
  Users,
} from "lucide-react";

export const siteConfig = {
  description:
    "Beautifully designed components that you can copy and paste into your apps. Accessible. Customizable. Open Source.",

  links: {
    "/": [
      {
        href: "/#features",
        Icon: Compass,
        title: "Features",
      },
      {
        href: "/#testimonials",
        Icon: Award,
        title: "Testimonials",
      },
      {
        href: "/#pricing",
        Icon: DollarSign,
        title: "Pricing",
      },
      {
        badge: undefined,
        href: "/#faq",
        Icon: HelpCircle,
        title: "FAQ",
      },
    ],
    "/chatbots": [
      {
        href: "/chatbots",
        Icon: MessageSquare,
        section: "create",
        title: "Chatbot",
      },
    ],
    "/chatbots/single": [
      {
        href: "/chatbots/single/chat",
        Icon: MessageSquare,
        tab: "Chat",
        title: "Chat",
      },
      {
        href: "/chatbots/single/threads",
        Icon: Layers,
        tab: "Threads",
        title: "Threads",
      },
      {
        href: "/chatbots/single/configuration",
        Icon: Bolt,
        tab: "Configuration",
        title: "Configuration",
      },
    ],
    "/content": [
      {
        href: "/content",
        Icon: Server,
        section: "data",
        title: "Content",
      },
    ],
    "/content/single": [
      {
        href: "/content/single/general",
        Icon: Server,
        tab: "General",
        title: "General",
      },
      {
        href: "/content/single/vectors",
        Icon: Server,
        tab: "Vectors",
        title: "Vectors",
      },
    ],
    "/home": [
      {
        href: "/home",
        Icon: HomeIcon,
        section: "home",
        title: "Home",
      },
    ],
    "/images": [
      {
        href: "/images",
        Icon: Image,
        section: "create",
        title: "Image",
      },
    ],
    "/import": [
      {
        href: "/import/file",
        Icon: CloudUploadIcon,
        section: "data",
        tab: "From File",
        title: "Import",
      },
      {
        href: "/import/url",
        Icon: Compass,
        tab: "From URL",
        title: "From URL",
      },
    ],
    "/settings/organization": [
      {
        href: "/settings/organization/general",
        Icon: Building2,
        section: "settings",
        tab: "General",
        title: "Organization",
      },
      {
        href: "/settings/organization/billing",
        Icon: DollarSign,
        tab: "Billing",
        title: "Billing",
      },
      {
        href: "/settings/organization/members",
        Icon: Users,
        tab: "Members",
        title: "Members",
      },
      {
        href: "/settings/organization/api-tokens",
        Icon: Key,
        tab: "API Tokens",
        title: "API Tokens",
      },
    ],
    "/settings/profile": [
      {
        href: "/settings/profile/general",
        Icon: User,
        section: "settings",
        tab: "General",
        title: "Profile",
      },
      {
        href: "/settings/profile/security",
        Icon: KeySquare,
        tab: "Security",
        title: "Security",
      },
    ],
  } as Record<
    string,
    {
      href: string;
      Icon: any;
      section?: string;
      tab?: string;
      title: string;
    }[]
  >,
  name: "Arches AI",
  ogImage: "https://ui.shadcn.com/og.jpg",
  url: "https://archesai.com",
};

export type SiteConfig = typeof siteConfig;
