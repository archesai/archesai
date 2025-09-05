/* eslint-disable @typescript-eslint/no-unsafe-assignment */

import { createElement, useCallback, useMemo } from "react";
import { Link as TanStackLink } from "@tanstack/react-router";

import type { SmartLinkProps } from "@archesai/ui/hooks/use-link";

import { filterUndefined, useLinkContext } from "@archesai/ui/hooks/use-link";

/**
 * Smart Link that automatically chooses between:
 * - TanStack Router Link for internal routes
 * - Regular anchor tags for external links
 * - Supports analytics tracking
 * - Handles loading states
 */
export const SmartLink: React.FC<SmartLinkProps> = ({
  children,
  className,
  disabled,
  download,
  // External link props
  external,
  hash,
  href,
  mask,
  newTab = true,
  noOpener = true,
  noReferrer = false,
  // Standard props
  onClick,
  params,
  preload = "intent",
  preloadDelay,
  preserveSearch,
  replace,
  resetScroll,
  resetSearch,
  // TanStack Router props
  search,
  state,
  // Analytics props
  trackClick = true,
  trackingAction = "click",
  trackingCategory = "Navigation",
  trackingLabel,
  trackingValue,
  ...restProps
}) => {
  const { isExternalUrl, trackEvent } = useLinkContext();
  // const routerState = useRouterState()

  // Determine if this is an external link
  const isExternal = external ?? isExternalUrl(href);

  // Check if we're currently navigating to this route
  // const isNavigating =
  //   routerState.status === 'pending' && routerState.location.pathname === href

  // Handle analytics tracking
  const handleClick = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      if (trackClick) {
        trackEvent(
          trackingCategory,
          trackingAction,
          trackingLabel ?? href,
          trackingValue,
        );
      }

      onClick?.(e);
    },
    [
      trackClick,
      trackingCategory,
      trackingAction,
      trackingLabel,
      href,
      trackingValue,
      trackEvent,
      onClick,
    ],
  );

  // Handle search params
  const linkSearch = useMemo<
    boolean | Record<string, unknown> | undefined
  >(() => {
    if (preserveSearch && search === undefined) {
      return true;
    }
    if (resetSearch) {
      return {};
    }
    return search as Record<string, unknown> | undefined;
  }, [search, preserveSearch, resetSearch]);

  // Build CSS classes
  const combinedClassName = useMemo(() => {
    const classes = [className];

    // if (isNavigating) {
    //   classes.push('opacity-50', 'pointer-events-none')
    // }

    if (disabled) {
      classes.push("opacity-50", "cursor-not-allowed", "pointer-events-none");
    }

    return classes.filter(Boolean).join(" ") || undefined;
  }, [
    className,
    // isNavigating,
    disabled,
  ]);

  // Filter out undefined values from rest props
  const safeRestProps = useMemo(() => filterUndefined(restProps), [restProps]);

  // Build TanStack Router props safely
  const tanstackProps = useMemo(() => {
    const props: Record<string, unknown> = { to: href };

    // Only add properties if they're not undefined
    if (linkSearch !== undefined) props.search = linkSearch;
    if (params !== undefined) props.params = params;
    if (hash !== undefined) props.hash = hash;
    if (state !== undefined) props.state = state;
    if (mask !== undefined) props.mask = mask;
    if (replace !== undefined) props.replace = replace;
    if (resetScroll !== undefined) props.resetScroll = resetScroll;
    props.preload = preload;
    if (preloadDelay !== undefined) props.preloadDelay = preloadDelay;
    if (disabled !== undefined) props.disabled = disabled;
    if (combinedClassName !== undefined) props.className = combinedClassName;
    props.onClick = handleClick;

    // Add safe rest props
    Object.assign(props, safeRestProps);

    return props;
  }, [
    href,
    linkSearch,
    params,
    hash,
    state,
    mask,
    replace,
    resetScroll,
    preload,
    preloadDelay,
    disabled,
    combinedClassName,
    handleClick,
    safeRestProps,
  ]);

  // External link rendering
  if (isExternal) {
    const externalProps = filterUndefined({
      className: combinedClassName,
      download: download === true ? "" : download,
      href,
      onClick: handleClick,
      rel:
        [noOpener && "noopener", noReferrer && "noreferrer"]
          .filter(Boolean)
          .join(" ") || undefined,
      target: newTab ? "_blank" : undefined,
      ...safeRestProps,
    });

    return createElement("a", externalProps, children);
  }

  // Internal link with TanStack Router
  return createElement(
    TanStackLink,
    tanstackProps,
    // isNavigating ?
    //   createElement(
    //     'span',
    //     { className: 'flex items-center gap-2' },
    //     createElement('span', {
    //       className:
    //         'animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full'
    //     }),
    //     children
    //   )
    // : children
    children,
  );
};
