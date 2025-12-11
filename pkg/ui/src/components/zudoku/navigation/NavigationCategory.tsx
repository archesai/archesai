import { ChevronRightIcon } from "lucide-react";
import { Collapsible as CollapsiblePrimitive } from "radix-ui";
import { memo, useEffect, useState } from "react";
import { Button } from "#components/shadcn/button";

const Collapsible = {
  Content: CollapsiblePrimitive.Content,
  Root: CollapsiblePrimitive.Root,
  Trigger: CollapsiblePrimitive.Trigger,
};

import type { NavigationItem as NavigationCategoryType } from "../context/ZudokuContext";
import { NavLink } from "../Link";
import { cn, joinUrl } from "../utils";
import { NavigationItem } from "./NavigationItem";
import { navigationListItem, useIsCategoryOpen } from "./utils";

const NavigationCategoryInner = ({
  category,
  onRequestClose,
}: {
  category: NavigationCategoryType;
  onRequestClose?: () => void;
}) => {
  const isCategoryOpen = useIsCategoryOpen(category);
  const [hasInteracted, setHasInteracted] = useState(false);
  const location = {
    pathname: window.location.pathname,
    search: window.location.search,
  };

  const isCollapsible = !(
    "collapsible" in category && category.collapsible === false
  );
  const isCollapsed = !(
    "collapsed" in category && category.collapsed === false
  );
  const isDefaultOpen = Boolean(
    !isCollapsible || !isCollapsed || isCategoryOpen,
  );
  const [open, setOpen] = useState(isDefaultOpen);
  const isActive = category.path && location.pathname === category.path;

  useEffect(() => {
    // this is triggered when an item from the navigation is clicked
    // and the navigation, enclosing this item, is not opened
    if (isCategoryOpen) {
      setOpen(true);
    }
  }, [isCategoryOpen]);

  const ToggleButton = isCollapsible && (
    <Button
      className="size-6 hover:bg-[hsl(from_var(--accent)_h_s_calc(l+6*var(--dark)))]"
      onClick={(e) => {
        e.preventDefault();
        setOpen((prev) => !prev);
        setHasInteracted(true);
      }}
      size="icon"
      variant="ghost"
    >
      <ChevronRightIcon
        className={cn(
          hasInteracted && "transition",
          "shrink-0 group-data-[state=open]:rotate-90 rtl:rotate-180",
        )}
        size={16}
      />
    </Button>
  );

  const icon = category.icon && (
    <span className={cn("align-[-0.125em]", isActive && "text-primary")}>
      {category.icon}
    </span>
  );

  const styles = navigationListItem({
    className: [
      "group text-start font-medium",
      isCollapsible || category.path !== undefined
        ? "cursor-pointer"
        : "cursor-default hover:bg-transparent",
    ],
  });

  return (
    <Collapsible.Root
      className="flex flex-col"
      defaultOpen={isDefaultOpen}
      onOpenChange={() => setOpen(true)}
      open={open}
    >
      <Collapsible.Trigger
        asChild
        className="group"
        disabled={!isCollapsible}
      >
        {category.path ? (
          <NavLink
            className={styles}
            onClick={() => {
              setHasInteracted(true);
              // if it is the current path and closed then open it because there's no path change to trigger the open
              if (isActive && !open) {
                setOpen(true);
              }
            }}
            to={joinUrl(category.path) + location.search}
          >
            {icon}
            <div className="flex w-full items-center justify-between gap-2 text-foreground/80 group-aria-[current='page']:text-primary">
              <div className="truncate">{category.label}</div>
              {ToggleButton}
            </div>
          </NavLink>
        ) : (
          // biome-ignore lint/a11y/noStaticElementInteractions: This is only to track if the user has interacted
          <div
            className={styles}
            onClick={() => setHasInteracted(true)}
            onKeyUp={(e) => {
              if (e.key === "Enter" || e.key === " ") setHasInteracted(true);
            }}
          >
            {icon}
            <div className="flex w-full items-center justify-between">
              <div className="flex w-full gap-2 truncate">{category.label}</div>
              {ToggleButton}
            </div>
          </div>
        )}
      </Collapsible.Trigger>
      <Collapsible.Content
        className={cn(
          // CollapsibleContent class is used to animate and it should only be applied when the user has triggered the toggle
          hasInteracted && "CollapsibleContent",
          (!category.children || category.children.length === 0) && "hidden",
          "my-1 ms-6",
        )}
      >
        <ul className="relative after:absolute after:-start-(--padding-nav-item) after:top-0 after:bottom-0 after:w-px after:translate-x-[1.5px] after:bg-border">
          {(category.children || []).map((item) => (
            <NavigationItem
              item={item}
              key={
                item.type +
                (item.label ?? "") +
                (item.path || "") +
                (item.file || "") +
                (item.to || "")
              }
              {...(onRequestClose && { onRequestClose })}
            />
          ))}
        </ul>
      </Collapsible.Content>
    </Collapsible.Root>
  );
};

export const NavigationCategory = memo(NavigationCategoryInner);

NavigationCategory.displayName = "NavigationCategory";
