import {
  Award,
  // Bolt,
  Building2,
  Compass,
  CreditCard,
  DollarSign,
  HelpCircle,
  HomeIcon,
  Image,
  KeySquare,
  // Layers,
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
        href: "/chatbots/single",
        Icon: MessageSquare,
        title: "Chat",
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
    "/organization": [
      {
        href: "/organization",
        Icon: Building2,
        section: "settings",
        title: "Organization",
      },
      {
        href: "/organization/general",
        Icon: Building2,
        tab: "General",
        title: "Organization",
      },
      {
        href: "/organization/billing",
        Icon: CreditCard,
        tab: "Billing",
        title: "Billing",
      },
      {
        href: "/organization/members",
        Icon: Users,
        tab: "Members",
        title: "Members",
      },
      {
        href: "/organization/api-tokens",
        Icon: Lock,
        tab: "API Tokens",
        title: "API Tokens",
      },
    ],
    "/profile": [
      {
        href: "/profile",
        Icon: User,
        section: "settings",
        title: "Profile",
      },
      {
        href: "/profile/general",
        Icon: User,
        tab: "General",
        title: "Profile",
      },
      {
        href: "/profile/security",
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
