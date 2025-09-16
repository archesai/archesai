import type { LinkProps as UILinkProps } from "@archesai/ui";
import { Link as UILink } from "@archesai/ui";
import type { LinkProps as RouterLinkProps } from "@tanstack/react-router";
import { Link as RouterLink } from "@tanstack/react-router";

type PlatformLinkProps = Omit<UILinkProps, "href"> & {
  to?: RouterLinkProps["to"];
  params?: RouterLinkProps["params"];
  search?: RouterLinkProps["search"];
  hash?: RouterLinkProps["hash"];
  state?: RouterLinkProps["state"];
  replace?: RouterLinkProps["replace"];
  preload?: RouterLinkProps["preload"];
};

/**
 * Platform-specific Link component that combines TanStack Router functionality
 * with the UI library's Link styling
 */
export function Link({
  to,
  params,
  search,
  hash,
  state,
  replace,
  preload,
  children,
  className,
  isActive: _isActive,
  external,
}: PlatformLinkProps) {
  // For external links, use the UI Link directly
  if (
    external ||
    (typeof to === "string" &&
      (to.startsWith("http://") || to.startsWith("https://")))
  ) {
    return (
      <UILink
        className={className}
        external={true}
        href={to as string}
      >
        {children}
      </UILink>
    );
  }

  // For internal links, use TanStack Router with explicit spreading
  // to avoid TypeScript exactOptionalPropertyTypes issues
  return (
    <RouterLink
      className={className}
      to={to || "/"}
      {...(hash !== undefined && { hash })}
      {...(params !== undefined && { params })}
      {...(preload !== undefined && { preload })}
      {...(replace !== undefined && { replace })}
      {...(search !== undefined && { search })}
      {...(state !== undefined && { state })}
    >
      {children}
    </RouterLink>
  );
}
