import type { ZudokuConfig } from "zudoku";
import { Button } from "zudoku/ui/Button";
import { docs } from "#sidebar";

const config: ZudokuConfig = {
  apis: [
    {
      input: "./apis/openapi.yaml",
      options: {
        expandAllTags: true,
      },
      path: "/docs/api",
      type: "file",
    },
  ],
  canonicalUrlOrigin: "https://archesai.github.io",
  docs: {
    defaultOptions: {
      disablePager: false,
      showLastModified: true,
      suggestEdit: {
        text: "Edit this page",
        url: "https://github.com/archesai/archesai/edit/main/docs/{filePath}",
      },
      toc: true,
    },
    files: ["/pages/**/*.{md,mdx}"],
  },
  enableStatusPages: true,
  metadata: {
    applicationName: "Arches AI",
    defaultTitle: "Arches AI",
    favicon: "/favicon.ico",
    logo: "/large-logo.svg",
    title: "%s | Arches AI",
  },
  navigation: docs,
  port: 3002,
  // search: {
  //   type: "pagefind",
  // },
  site: {
    banner: {
      color: "#7362ef",
      dismissible: true,
      message: "⭐️ If you like Arches, give it a star on GitHub! ⭐️",
    },
    footer: {
      copyright: `© ${new Date().getFullYear()} Arches AI`,
      logo: {
        alt: "Arches AI",
        href: "/getting-started",
        src: {
          dark: "/large-logo-white.svg",
          light: "/large-logo.svg",
        },
      },
      social: [
        {
          href: "github.com/archesai/archesai",
          label: "GitHub",
        },
        {
          href: "twitter.com/archesai",
          label: "Twitter",
        },
      ],
    },
    logo: {
      alt: "Arches AI",
      href: "/getting-started",
      src: {
        dark: "/large-logo-white.svg",
        light: "/large-logo.svg",
      },
      width: "130px",
    },
    logoUrl: "/large-logo.svg",
    showPoweredBy: false,
    title: "Arches AI",
  },
  sitemap: {
    siteUrl: "https://archesai.github.io",
  },

  slots: {
    "head-navigation-end": () => (
      <div className="flex items-center border-r pe-2">
        <Button
          asChild
          size="icon"
          variant="ghost"
        >
          <a
            aria-label="Visit Arches on GitHub"
            href="https://github.com/archesai/archesai"
            rel="noopener noreferrer"
          >
            {/* <GithubIcon
              aria-hidden="true"
              className="w-4 h-4 dark:invert"
            /> */}
          </a>
        </Button>
        <Button
          asChild
          size="icon"
          variant="ghost"
        >
          <a
            aria-label="Join Arches Discord community"
            href="https://discord.archesai.dev"
            rel="noopener noreferrer"
          >
            {/* <DiscordIcon
              aria-hidden="true"
              className="w-5 h-5 dark:invert"
            /> */}
          </a>
        </Button>
      </div>
    ),
  },

  theme: {
    customCss: `

@theme {
  --animate-wiggle: wiggle 1s ease-in-out infinite;
  @keyframes wiggle {
    0%,
    100% {
      transform: rotate(-3deg);
    }
    50% {
      transform: rotate(3deg);
    }
  }
}`,

    dark: {
      primary: "#7362ef",
      primaryForeground: "#000000",
    },
    fonts: {
      mono: "Geist Mono",
      sans: "Geist",
      serif: "Geist",
    },
    light: {
      primary: "#7362ef",
      primaryForeground: "#FFFFFF",
    },
  },
};

export default config;
