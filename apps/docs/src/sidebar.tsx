/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import type { Navigation } from "zudoku";
import { LandingPage } from "#landing";

/**
 * NOTE: This file should not import anything except zudoku. We use this file
 * in the build of the zuplo docs site to generate the sidebar there.
 */

export const docs: Navigation = [
  {
    display: "hide",
    element: <LandingPage />,
    label: "Landing Page",
    layout: "none",
    path: "/",
    type: "custom-page",
  },
  {
    items: [
      {
        icon: "sparkles",
        items: [
          "/getting-started",
          "/guides/quickstart",
          "/cli-reference",
          "/contributing",
        ],
        label: "Getting Started",
        type: "category",
      },
      {
        icon: "layers",
        items: [
          "/architecture/overview",
          "/architecture/system-design",
          "/architecture/code-generation",
          "/architecture/authentication",
          "/architecture/project-layout",
        ],
        label: "Architecture",
        type: "category",
      },
      {
        icon: "package",
        items: [
          "/features/overview",
          "/features/auth",
          "/features/organizations",
          "/features/workflows",
          "/features/content",
          "/features/tui",
        ],
        label: "Features",
        type: "category",
      },
      {
        icon: "book-open",
        items: [
          "/guides/overview",
          "/guides/development",
          "/guides/code-generation",
          "/guides/testing",
          "/guides/custom-handlers",
          "/guides/makefile-commands",
        ],
        label: "Guides",
        type: "category",
      },
      {
        icon: "rocket",
        items: [
          "/deployment/overview",
          "/deployment/docker",
          "/deployment/kubernetes",
          "/deployment/production",
        ],
        label: "Deployment",
        type: "category",
      },
      {
        icon: "shield",
        items: ["/security/overview", "/security/best-practices"],
        label: "Security",
        type: "category",
      },
      {
        icon: "wrench",
        items: ["/troubleshooting/common-issues"],
        label: "Troubleshooting",
        type: "category",
      },
      {
        icon: "map",
        items: ["/ROADMAP"],
        label: "Roadmap",
        type: "category",
      },
      {
        icon: "code",
        items: [
          "/api-reference/overview",
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
  {
    // element: <ThemePlayground />,
    element: <></>,
    label: "Themes",
    path: "/docs/theme-playground",
    type: "custom-page",
  },
];
