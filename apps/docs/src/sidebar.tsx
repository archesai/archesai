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
          "/configuration",
          "/contributing",
        ],
        label: "Getting Started",
        type: "category",
      },
      {
        icon: "package",
        items: ["/features/authentication", "/features/database"],
        label: "Features",
        type: "category",
      },
      {
        icon: "book-open",
        items: ["/guides/code-generation", "/guides/custom-handlers"],
        label: "Guides",
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
