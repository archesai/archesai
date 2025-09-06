import type { SiteRoute } from "@archesai/ui/lib/site-config.interface";

import { siteRoutes } from "#lib/site-config";

export const flattenRoutes = (routes: SiteRoute[]): SiteRoute[] => {
  return routes.flatMap((route) =>
    route.children ? [route, ...route.children] : [route],
  );
};

export const getRouteMeta = (
  pathname: string,
): {
  description: string;
  Icon: React.ComponentType<Record<string, unknown>> | undefined;
  title: string;
} => {
  const allRoutes = flattenRoutes(siteRoutes);
  const matched = allRoutes.find((route) => route.href === pathname);

  return {
    description: matched?.description ?? "",
    Icon: matched?.Icon,
    title: matched?.title ?? "",
  };
};

export const getTabsForPath = (pathname: string): SiteRoute[] => {
  const matched = siteRoutes.find((route) => pathname.startsWith(route.href));

  return matched?.children?.filter((r) => r.showInTabs) ?? [];
};

export const getActiveTab = (pathname: string): string | undefined => {
  const tabs = getTabsForPath(pathname);
  return tabs.find((tab) => tab.href === pathname)?.href;
};

export const getSections = (): string[] =>
  Array.from(new Set(siteRoutes.map((r) => r.section)));
