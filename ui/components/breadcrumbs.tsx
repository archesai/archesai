import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { siteConfig } from "@/config/site";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react"; // Adjust the import path as necessary

import { getOriginalPath } from "./tabs-section";

export const Breadcrumbs = () => {
  const pathname = usePathname() as string;
  const pathnames = pathname.split("/").filter((x) => x);

  let Icon = null;
  const entries = Object.entries(siteConfig.links);
  entries.sort(([a], [b]) => a.length - b.length);
  for (const [routePattern, links] of entries) {
    if (
      pathname.startsWith(routePattern) ||
      getOriginalPath(pathname).startsWith(routePattern)
    ) {
      const link = links.find(
        (link) =>
          link.href === pathname || link.href === getOriginalPath(pathname)
      );
      Icon = link?.Icon;
    }
  }

  return (
    <Breadcrumb>
      <BreadcrumbList>
        <Icon className="h-5 w-5  text-muted-foreground" strokeWidth={1.5} />
        {pathnames.map((value, index) => {
          const to = `/${pathnames.slice(0, index + 1).join("/")}`;
          // const isLast = index === pathnames.length - 1;
          const isLast = false;

          const capitalizedValue = (
            value.charAt(0).toUpperCase() + value.slice(1)
          ).slice(0, 13);
          return (
            <React.Fragment key={to}>
              <BreadcrumbItem>
                {isLast ? (
                  <Link className="text-foreground" href={to}>
                    {capitalizedValue}
                  </Link>
                ) : (
                  <Link href={to}>{capitalizedValue}</Link>
                )}
              </BreadcrumbItem>
              {index < pathnames.length - 1 && <BreadcrumbSeparator />}
            </React.Fragment>
          );
        })}
      </BreadcrumbList>
    </Breadcrumb>
  );
};
