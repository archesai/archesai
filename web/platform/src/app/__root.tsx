/// <reference types="vite/client" />

import type { GetSession200 } from "@archesai/client";
import { getGetSessionQueryKey } from "@archesai/client";
import { ThemeProvider, Toaster } from "@archesai/ui";
import { seo } from "@archesai/ui/lib/seo";
import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts,
} from "@tanstack/react-router";
import type { JSX } from "react";
import { DefaultCatchBoundary } from "#components/default-catch-boundary";
import NotFound from "#components/not-found";
import getServerSession from "#lib/get-headers";
import globalsCss from "../styles/globals.css?url";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  session: GetSession200 | null;
}>()({
  beforeLoad: async ({ context }) => {
    const session = await context.queryClient.fetchQuery({
      queryFn: ({ signal }) =>
        getServerSession({
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
        description:
          "Arches AI is the perfect tool to explore documents using artificial intelligence. Simply upload your PDF and start asking questions to your personalized chatbot.",
        image: "https://www.archesai.com/sc.png",
        title: "Arches AI",
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
