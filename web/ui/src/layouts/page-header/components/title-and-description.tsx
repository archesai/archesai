import type { JSX } from "react";

import type { SiteRoute } from "#lib/site-config.interface";

export const TitleAndDescription = ({
  siteRoute,
}: {
  siteRoute: SiteRoute;
}): JSX.Element | null => {
  if (!siteRoute.title) return null;
  return (
    <div className="container flex items-center gap-3 border-b px-4 py-3">
      <siteRoute.Icon className="h-8 w-8" />
      <div>
        <p className="font-semibold text-foreground/85 text-xl">
          {siteRoute.title}
        </p>
        <p className="text-muted-foreground text-sm">{siteRoute.description}</p>
      </div>
    </div>
  );
};
