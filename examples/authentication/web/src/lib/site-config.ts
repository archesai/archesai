import {
  BookOpenIcon,
  Building2Icon,
  CpuIcon,
  FileIcon,
  PackageCheckIcon,
  ServerIcon,
  Settings2Icon,
  SparklesIcon,
  SquareTerminalIcon,
  TagsIcon,
  UserIcon,
  UsersIcon,
  WorkflowIcon,
} from "@archesai/ui";
import type { SiteRoute } from "@archesai/ui/lib/site-config.interface";

export const siteRoutes: SiteRoute[] = [
  {
    description: "Try out your tools here.",
    href: "/",
    Icon: SquareTerminalIcon,
    section: "Home",
    title: "Dashboard",
  },
  {
    description: "Browse and manage your artifacts here.",
    href: "/artifacts",
    Icon: ServerIcon,
    section: "Data",
    title: "Artifacts",
  },
  {
    description: "Explore and run tools.",
    href: "/tools",
    Icon: CpuIcon,
    section: "Build",
    title: "Tools",
  },
  {
    description: "View and manage your pipelines.",
    href: "/pipelines",
    Icon: WorkflowIcon,
    section: "Build",
    showInTabs: true,
    title: "Pipelines",
  },
  {
    description: "View your previous runs.",
    href: "/runs",
    Icon: PackageCheckIcon,
    section: "Build",
    title: "History",
  },
  {
    description: "Create and manage labels.",
    href: "/labels",
    Icon: TagsIcon,
    section: "Data",
    title: "Labels",
  },
  {
    description: "View and explore the OpenAPI specification.",
    href: "/openapi",
    Icon: BookOpenIcon,
    section: "Build",
    title: "API Docs",
  },

  {
    description: "View your profile information.",
    href: "/profile",
    Icon: UserIcon,
    section: "Settings",
    title: "Account",
  },
  {
    description: "View and explore ArchesAI configuration schema.",
    href: "/configuration",
    Icon: FileIcon,
    section: "Settings",
    title: "Configuration",
  },
  {
    children: [
      {
        description: "View and update your organization's general information.",
        href: "/organization",
        Icon: Building2Icon,
        section: "Settings",
        showInTabs: true,
        title: "General",
      },
      {
        description: "View and manage your organization's members.",
        href: "/organization/members",
        Icon: UsersIcon,
        section: "Settings",
        showInTabs: true,
        title: "Members",
      },
      {
        description: "",
        href: "/profile/themes",
        Icon: SparklesIcon,
        section: "Settings",
        showInTabs: true,
        title: "Themes",
      },
    ],
    href: "/organization",
    Icon: Settings2Icon,
    section: "Settings",
    title: "Settings",
  },
];

export const siteMetadata = {
  description:
    "Arches AI is a platform that provides tools to transform data into various forms of content.",
  name: "Arches AI",
  ogImage: "https://ui.shadcn.com/og.jpg",
  url: "https://archesai.com",
};
