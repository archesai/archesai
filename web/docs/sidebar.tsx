import type { Navigation } from "zudoku";

/**
 * NOTE: This file should not import anything except zudoku. We use this file
 * in the build of the zuplo docs site to generate the sidebar there.
 */

export const docs: Navigation = [
  "docs/quickstart",
  // {
  //   type: "category",
  //   label: "Getting started",
  //   icon: "circle-play",
  //   items: [
  //     {
  //       type: "custom-page",
  //       path: "docs/introduction",
  //       label: "Introduction",
  //       element: <Introduction />,
  //     },
  //   ],
  // },
  {
    icon: "settings",
    items: [
      "docs/customization/colors-theme",
      "docs/configuration/docs",
      "docs/configuration/navigation",
      "docs/configuration/site",
      "docs/configuration/search",
      "docs/configuration/footer",
    ],
    label: "Configuration",
    link: "docs/configuration/overview",
    type: "category",
  },
  {
    icon: "book-open-text",
    items: [
      "docs/markdown/overview",
      "docs/markdown/frontmatter",
      "docs/markdown/mdx",
      "docs/markdown/admonitions",
      "docs/markdown/code-blocks",
    ],
    label: "Writing",
    link: "docs/writing",
    type: "category",
  },
  {
    icon: "shapes",
    items: ["docs/concepts/auth-provider-api-identities"],
    label: "Concepts",
    type: "category",
  },
  {
    icon: "globe",
    items: [
      "docs/configuration/api-reference",
      "docs/configuration/api-catalog",
    ],
    label: "OpenAPI",
    type: "category",
  },
  {
    icon: "lock",
    items: [
      "docs/configuration/authentication",
      "docs/configuration/authentication-auth0",
      "docs/configuration/authentication-clerk",
      "docs/configuration/authentication-azure-ad",
      "docs/configuration/authentication-pingfederate",
      "docs/configuration/authentication-supabase",
    ],
    label: "Authentication",
    type: "category",
  },
  {
    icon: "blocks",
    items: ["docs/configuration/sentry"],
    label: "Integrations",
    type: "category",
  },
  {
    icon: "monitor-check",
    items: [
      "docs/guides/static-files",
      "docs/guides/environment-variables",
      "docs/guides/custom-pages",
      "docs/guides/navigation-migration",
      "docs/guides/using-multiple-apis",
      "docs/guides/managing-api-keys-and-identities",
      "docs/guides/transforming-examples",
      "docs/guides/processors",
    ],
    label: "Guides",
    type: "category",
  },
  {
    icon: "cloud-upload",
    items: [
      "docs/deploy/cloudflare-pages",
      "docs/deploy/github-pages",
      "docs/deploy/vercel",
      "docs/deploy/direct-upload",
    ],
    label: "Deployment",
    link: "docs/deployment",
    type: "category",
  },
  {
    icon: "blocks",
    items: [
      "docs/configuration/build-configuration",
      "docs/configuration/vite-config",
      "docs/configuration/slots",
      "docs/custom-plugins",
      "docs/extending/events",
    ],
    label: "Extending",
    type: "category",
  },
];
export const components: Navigation = [
  {
    icon: "album",
    items: [
      "docs/components/alert",
      "docs/components/callout",
      "docs/components/badge",
      "docs/components/card",
      "docs/components/icons",
      "docs/components/markdown",
      "docs/components/typography",
    ],
    label: "General",
    type: "category",
  },
  {
    icon: "component",
    items: [
      "docs/components/playground",
      "docs/components/secret",
      "docs/components/stepper",
      "docs/components/syntax-highlight",
    ],
    label: "Documentation",
    type: "category",
  },
  {
    icon: "text-cursor-input",
    items: [
      "docs/components/button",
      "docs/components/checkbox",
      "docs/components/slider",
      "docs/components/switch",
      "docs/components/label",
    ],
    label: "Form",
    type: "category",
  },
  {
    icon: "hard-hat",
    items: [
      "docs/components/client-only",
      "docs/components/head",
      "docs/components/slot",
    ],
    label: "Utility",
    type: "category",
  },
];
