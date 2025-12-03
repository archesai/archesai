import { ExternalLinkIcon } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "#components/shadcn/tooltip";
import { AnchorLink } from "../anchor-link";
import { useViewportAnchor } from "../context/ViewportAnchorContext";
import type { NavigationItem as NavigationItemType } from "../context/ZudokuContext";
import { NavLink } from "../Link";
import { cn, joinUrl, shouldShowItem } from "../utils";
import { NavigationBadge } from "./NavigationBadge";

type BadgeType = {
  color:
    | "blue"
    | "gray"
    | "green"
    | "indigo"
    | "outline"
    | "purple"
    | "red"
    | "yellow";
  label: string;
  className?: string;
  invert?: boolean;
};

const hasBadge = (item: unknown): item is { badge: BadgeType } => {
  if (typeof item !== "object" || item === null) {
    return false;
  }

  if (!("badge" in item)) {
    return false;
  }

  const itemWithBadge = item as { badge: unknown };
  const badge = itemWithBadge.badge;

  if (typeof badge !== "object" || badge === null) {
    return false;
  }

  if (!("color" in badge) || !("label" in badge)) {
    return false;
  }

  const badgeWithProps = badge as { label: unknown; color: unknown };
  return typeof badgeWithProps.label === "string";
};

import { NavigationCategory } from "./NavigationCategory";
import { navigationListItem } from "./utils";

const TruncatedLabel = ({
  label,
  className,
}: {
  label: string;
  className?: string;
}) => {
  const ref = useRef<HTMLSpanElement>(null);
  const [isTruncated, setIsTruncated] = useState(false);

  useEffect(() => {
    if (!ref.current) return;

    if (ref.current.offsetWidth < ref.current.scrollWidth) {
      setIsTruncated(true);
    }
  }, []);

  return (
    <>
      <span
        className={cn("flex-1 truncate", className)}
        ref={ref}
        title={label}
      >
        {label}
      </span>
      {isTruncated && (
        <TooltipProvider delayDuration={500}>
          <Tooltip disableHoverableContent>
            <TooltipTrigger className="absolute inset-0 z-10" />
            <TooltipContent
              align="center"
              className="max-w-64 rounded-lg"
              side="bottom"
            >
              {label}
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}
    </>
  );
};

export const DATA_ANCHOR_ATTR = "data-anchor";

export const NavigationItem = ({
  item,
  onRequestClose,
}: {
  item: NavigationItemType;
  onRequestClose?: () => void;
}) => {
  const location = { pathname: window.location.pathname };
  const { activeAnchor } = useViewportAnchor();

  if (!shouldShowItem(item, true)) {
    return null;
  }

  switch (item.type) {
    case "category":
      return onRequestClose ? (
        <NavigationCategory
          category={item}
          onRequestClose={onRequestClose}
        />
      ) : (
        <NavigationCategory category={item} />
      );
    case "doc": {
      const isActive = location.pathname === joinUrl(item.path || "");
      return (
        <NavLink
          className={navigationListItem({ isActive })}
          {...(onRequestClose && { onClick: onRequestClose })}
          to={joinUrl(item.path || "")}
        >
          {item.icon && <span className="align-[-0.125em]">{item.icon}</span>}
          {hasBadge(item) ? (
            <>
              {item.label && (
                <TruncatedLabel
                  className="flex-1"
                  label={item.label}
                />
              )}
              <NavigationBadge {...(item as { badge: BadgeType }).badge} />
            </>
          ) : (
            item.label
          )}
        </NavLink>
      );
    }
    case "link":
    case "custom-page": {
      const href =
        item.type === "link" ? item.to || "" : joinUrl(item.path || "");
      const isActiveLink =
        href === [location.pathname, activeAnchor].filter(Boolean).join("#");
      return !href?.startsWith("http") ? (
        <AnchorLink
          to={href}
          {...{ [DATA_ANCHOR_ATTR]: href.split("#")[1] }}
          className={navigationListItem({
            isActive: isActiveLink,
          })}
          {...(onRequestClose && { onClick: onRequestClose })}
        >
          {item.icon && <span className="align-[-0.125em]">{item.icon}</span>}
          {hasBadge(item) ? (
            <>
              {item.label && <TruncatedLabel label={item.label} />}
              <NavigationBadge {...(item as { badge: BadgeType }).badge} />
            </>
          ) : (
            <span className="break-all">{item.label}</span>
          )}
        </AnchorLink>
      ) : (
        <a
          className={navigationListItem()}
          href={href}
          onClick={onRequestClose}
          rel="noopener noreferrer"
          target="_blank"
        >
          {item.icon && <span className="align-[-0.125em]">{item.icon}</span>}
          <span className="whitespace-normal">{item.label}</span>
          {/* This prevents that the icon would be positioned in its own line if the text fills a line entirely */}
          <span className="whitespace-nowrap">
            <ExternalLinkIcon
              className="-translate-y-0.5 inline"
              size={12}
            />
          </span>
        </a>
      );
    }
    default:
      return null;
  }
};
