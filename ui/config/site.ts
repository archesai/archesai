// config/site.ts
import {
  Award,
  Bolt,
  Building2,
  CloudUploadIcon,
  Compass,
  CreditCard,
  DollarSign,
  HelpCircle,
  HomeIcon,
  Image,
  KeySquare,
  Layers,
  Lock,
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
        href: "/content/single",
        Icon: Server,
        title: "Detils",
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
        href: "/import",
        Icon: CloudUploadIcon,
        section: "data",
        title: "Import",
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
        Icon: CreditCard,
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
        Icon: Lock,
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
