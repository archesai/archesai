import type { Navigation } from "zudoku";

/**
 * NOTE: This file should not import anything except zudoku. We use this file
 * in the build of the zuplo docs site to generate the sidebar there.
 */

export const docs: Navigation = [
  {
    items: [
      {
        icon: "sparkles",
        items: [
          "/documentation/getting-started",
          "/documentation/contributing",
        ],
        label: "Getting Started",
        type: "category",
      },
      {
        icon: "layers",
        items: [
          "/documentation/architecture/overview",
          "/documentation/architecture/system-design",
          "/documentation/architecture/project-layout",
          "/documentation/architecture/authentication",
        ],
        label: "Architecture",
        type: "category",
      },
      {
        icon: "package",
        items: [
          "/documentation/features/overview",
          "/documentation/features/auth",
          "/documentation/features/organizations",
          "/documentation/features/workflows",
          "/documentation/features/content",
          "/documentation/features/tui",
        ],
        label: "Features",
        type: "category",
      },
      {
        icon: "book-open",
        items: [
          "/documentation/guides/overview",
          "/documentation/guides/development",
          "/documentation/guides/testing",
          "/documentation/guides/test-coverage-report",
          "/documentation/guides/makefile-commands",
        ],
        label: "Guides",
        type: "category",
      },
      {
        icon: "rocket",
        items: [
          "/documentation/deployment/overview",
          "/documentation/deployment/docker",
          "/documentation/deployment/kubernetes",
          "/documentation/deployment/production",
          "/documentation/deployment/site",
        ],
        label: "Deployment",
        type: "category",
      },
      {
        icon: "shield",
        items: [
          "/documentation/security/overview",
          "/documentation/security/best-practices",
        ],
        label: "Security",
        type: "category",
      },
      {
        icon: "zap",
        items: [
          "/documentation/performance/overview",
          "/documentation/performance/optimization",
        ],
        label: "Performance",
        type: "category",
      },
      {
        icon: "activity",
        items: ["/documentation/monitoring/overview"],
        label: "Monitoring",
        type: "category",
      },
      {
        icon: "wrench",
        items: ["/documentation/troubleshooting/common-issues"],
        label: "Troubleshooting",
        type: "category",
      },
      {
        icon: "code",
        items: [
          "/documentation/api-reference/overview",
          {
            badge: {
              color: "green",
              label: "Interactive",
            },
            icon: "external-link",
            label: "OpenAPI Explorer",
            to: "/docs/api",
            type: "link",
          },
        ],
        label: "API Reference",
        type: "category",
      },
      {
        collapsible: false,
        icon: "link",
        items: [
          {
            icon: "github",
            label: "GitHub Repository",
            to: "https://github.com/archesai/archesai",
            type: "link",
          },
          {
            icon: "book",
            label: "Arches AI Platform",
            to: "https://archesai.com",
            type: "link",
          },
        ],
        label: "External Links",
        type: "category",
      },
    ],
    label: "Documentation",
    type: "category",
  },
  {
    label: "API Reference",
    to: "/docs/api",
    type: "link",
  },
];
