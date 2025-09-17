import type { JSX } from "react";

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "#components/shadcn/breadcrumb";

export interface BreadcrumbData {
  title: string;
  path?: string;
}

export interface BreadCrumbsProps {
  items?: BreadcrumbData[];
  currentPath?: string;
  onNavigate?: (path: string) => void;
}

export const BreadCrumbs = ({
  items = [],
  currentPath = "",
  onNavigate,
}: BreadCrumbsProps): JSX.Element => {
  // If no items provided, generate from currentPath
  const breadcrumbItems =
    items.length > 0
      ? items
      : (() => {
          const pathSegments = currentPath.split("/").filter(Boolean);
          return pathSegments.map((segment, index) => {
            const path = `/${pathSegments.slice(0, index + 1).join("/")}`;
            const title = segment
              .split("-")
              .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
              .join(" ");

            return {
              path: index === pathSegments.length - 1 ? undefined : path,
              title,
            };
          });
        })();

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {breadcrumbItems.map((breadcrumb, index) => (
          <BreadcrumbItem
            className={index === 0 ? "hidden md:flex" : "flex"}
            key={breadcrumb.path || breadcrumb.title}
          >
            {index > 0 && <BreadcrumbSeparator />}
            {!breadcrumb.path ? (
              <BreadcrumbPage>{breadcrumb.title}</BreadcrumbPage>
            ) : (
              <BreadcrumbLink
                href={breadcrumb.path}
                onClick={(e) => {
                  if (onNavigate && breadcrumb.path) {
                    e.preventDefault();
                    onNavigate(breadcrumb.path);
                  }
                }}
              >
                {breadcrumb.title}
              </BreadcrumbLink>
            )}
          </BreadcrumbItem>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  );
};
