/// <reference types="vite/client" />

import {
  DefaultCatchBoundary,
  NotFound,
  ThemeProvider,
  Toaster,
} from "@archesai/ui";
import { seo } from "@archesai/ui/lib/seo";
import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts,
} from "@tanstack/react-router";
import type { JSX } from "react";
import getSessionSSR from "#lib/get-session-ssr";
import type { GetSessionQueryResult } from "#lib/index";
import { getGetSessionQueryKey } from "#lib/index";
import globalsCss from "../styles/globals.css?url";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  session: GetSessionQueryResult | null;
}>()({
  beforeLoad: async ({ context }) => {
    const session = await context.queryClient.fetchQuery({
      queryFn: ({ signal }) =>
        getSessionSSR({
          signal,
        }),
      queryKey: getGetSessionQueryKey(),
    });
    return {
      session,
    };
  },
  component: RootComponent,
  errorComponent: (props) => {
    return (
      <RootDocument>
        <DefaultCatchBoundary {...props} />
      </RootDocument>
    );
  },
  head: () => ({
    links: [
      {
        href: globalsCss,
        rel: "stylesheet",
      },
      {
        href: "/apple-touch-icon.png",
        rel: "apple-touch-icon",
        sizes: "180x180",
      },
      {
        href: "/favicon-32x32.png",
        rel: "icon",
        sizes: "32x32",
        type: "image/png",
      },
      {
        href: "/favicon-16x16.png",
        rel: "icon",
        sizes: "16x16",
        type: "image/png",
      },
      {
        color: "#fffff",
        href: "/site.webmanifest",
        rel: "manifest",
      },
      {
        href: "/favicon.ico",
        rel: "icon",
      },
    ],
    meta: [
      {
        charSet: "utf-8",
      },
      {
        content: "width=device-width, initial-scale=1",
        name: "viewport",
      },
      ...seo({
        description: "Generated with Arches",
        title: "Basic",
      }),
    ],
  }),
  notFoundComponent: () => <NotFound />,
});

function RootDocument({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <html
      lang="en"
      suppressHydrationWarning
    >
      <head>
        <HeadContent />
      </head>
      <body>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          disableTransitionOnChange
          enableColorScheme
          enableSystem
        >
          {children}
          <Toaster />
        </ThemeProvider>
        <Scripts />
      </body>
    </html>
  );
}

function RootComponent() {
  return (
    <RootDocument>
      <Outlet />
    </RootDocument>
  );
}
