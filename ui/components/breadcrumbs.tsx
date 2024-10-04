import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react"; // Adjust the import path as necessary

export const Breadcrumbs = () => {
  const pathname = usePathname();
  const pathnames = pathname.split("/").filter((x) => x);

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {pathnames.map((value, index) => {
          const to = `/${pathnames.slice(0, index + 1).join("/")}`;
          const isLast = index === pathnames.length - 1;

          const capitalizedValue = (
            value.charAt(0).toUpperCase() + value.slice(1)
          ).slice(0, 13);
          return (
            <React.Fragment key={to}>
              <BreadcrumbItem>
                {isLast && index > 0 ? (
                  <Link
                    // className="text-foreground"
                    href={to}
                  >
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
