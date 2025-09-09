import type { ZudokuConfig } from "zudoku";

const config: ZudokuConfig = {
  // redirects: [
  //   { from: "/docs", to: "/docs/quickstart" },
  //   { from: "/getting-started", to: "/docs/quickstart" },
  //   { from: "/app-quickstart", to: "/docs/quickstart" },
  //   { from: "/components", to: "/components/callout" },
  //   { from: "/configuration/page", to: "/configuration/site" },
  // ],
  apis: [
    {
      input: "./apis/openapi.yaml",
      path: "/docs/api",
      type: "file",
    },
  ],
  canonicalUrlOrigin: "https://docs.archesai.com",
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
  metadata: {
    defaultTitle: "Arches AI",
    favicon: "https://platform.archesai.com/icon.png",
    title: "%s | Arches AI",
  },
  navigation: [
    {
      items: [
        {
          icon: "sparkles",
          items: [
            "/introduction",
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
  ],
  search: {
    type: "pagefind",
  },
  site: {
    logo: {
      alt: "Arches AI",
      src: { dark: "/logo-dark.svg", light: "/logo-light.svg" },
      width: "130px",
    },
    showPoweredBy: false,
  },
  sitemap: {
    siteUrl: "https://docs.archesai.com",
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
    light: {
      primary: "#7362ef",
      primaryForeground: "#FFFFFF",
    },
  },

  //   slots: {
  //   "head-navigation-end": () => (
  //     <div className="flex items-center border-r pe-2">
  //       <Button variant="ghost" size="icon" asChild>
  //         <a
  //           href="https://github.com/zuplo/zudoku"
  //           aria-label="Visit Zudoku on GitHub"
  //           rel="noopener noreferrer"
  //         >
  //           <GithubIcon className="w-4 h-4 dark:invert" aria-hidden="true" />
  //         </a>
  //       </Button>
  //       <Button variant="ghost" size="icon" asChild>
  //         <a
  //           href="https://discord.zudoku.dev"
  //           aria-label="Join Zudoku Discord community"
  //           rel="noopener noreferrer"
  //         >
  //           <DiscordIcon className="w-5 h-5 dark:invert" aria-hidden="true" />
  //         </a>
  //       </Button>
  //     </div>
  //   ),
  // },
  //   plugins: [
  //   {
  //     getHead: () => {
  //       return (
  //         <script>
  //           {`!function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(".");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement("script")).type="text/javascript",p.crossOrigin="anonymous",p.async=!0,p.src=s.api_host.replace(".i.posthog.com","-assets.i.posthog.com")+"/static/array.js",(r=t.getElementsByTagName("script")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a="posthog",u.people=u.people||[],u.toString=function(t){var e="posthog";return"posthog"!==a&&(e+="."+a),t||(e+=" (stub)"),e},u.people.toString=function(){return u.toString(1)+".people (stub)"},o="init capture register register_once register_for_session unregister unregister_for_session getFeatureFlag getFeatureFlagPayload isFeatureEnabled reloadFeatureFlags updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures on onFeatureFlags onSessionId getSurveys getActiveMatchingSurveys renderSurvey canRenderSurvey getNextSurveyStep identify setPersonProperties group resetGroups setPersonPropertiesForFlags resetPersonPropertiesForFlags setGroupPropertiesForFlags resetGroupPropertiesForFlags reset get_distinct_id getGroups get_session_id get_session_replay_url alias set_config startSessionRecording stopSessionRecording sessionRecordingStarted captureException loadToolbar get_property getSessionProperty createPersonProfile opt_in_capturing opt_out_capturing has_opted_in_capturing has_opted_out_capturing clear_opt_in_out_capturing debug getPageViewId captureTraceFeedback captureTraceMetric".split(" "),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]);
  //   posthog.init('phc_l8rjm0vHBMwNdGeBRDrK8UIYjyVxZyBAtnYo2hS18OY', { api_host: 'https://us.i.posthog.com', person_profiles: 'identified_only', })`}
  //         </script>
  //       );
  //     },
  //   },
  // ],
};

export default config;
