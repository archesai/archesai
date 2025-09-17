import { cva } from "class-variance-authority";
import type { NavigationItem } from "../context/ZudokuContext";
import { useCurrentNavigation } from "../context/ZudokuContext";
import { joinUrl } from "../utils";

const useLocation = () => {
  return { pathname: window.location.pathname };
};

export type TraverseCallback<T> = (
  item: NavigationItem,
  parentCategories: NavigationItem[],
) => T | undefined;

export const traverseNavigation = <T>(
  navigation: NavigationItem[],
  callback: TraverseCallback<T>,
): T | undefined => {
  for (const item of navigation) {
    const result = traverseNavigationItem(item, callback);
    if (result !== undefined) return result;
  }
  return undefined;
};

export const traverseNavigationItem = <T>(
  item: NavigationItem,
  callback: TraverseCallback<T>,
  parentCategories: NavigationItem[] = [],
): T | undefined => {
  const result = callback(item, parentCategories);
  if (result !== undefined) return result;

  if (item.type === "category" && item.children) {
    for (const child of item.children) {
      const childResult = traverseNavigationItem(child, callback, [
        ...parentCategories,
        item,
      ]);
      if (childResult !== undefined) return childResult;
    }
  }
  return undefined;
};

export const useCurrentItem = () => {
  const location = useLocation();
  const { navigation } = useCurrentNavigation();

  return traverseNavigation(navigation, (item) => {
    if (
      item.type === "doc" &&
      item.path &&
      joinUrl(item.path) === location.pathname
    ) {
      return item;
    }
    return undefined;
  });
};

export const useIsCategoryOpen = (category: NavigationItem) => {
  const location = useLocation();

  return traverseNavigationItem(category, (item) => {
    switch (item.type) {
      case "category":
        if (!item.path) {
          return undefined;
        }
        return joinUrl(item.path) === location.pathname ? true : undefined;
      case "custom-page":
      case "doc":
        return joinUrl(item.path || "") === location.pathname
          ? true
          : undefined;
      default:
        return undefined;
    }
  });
};

export const usePrevNext = (): {
  prev?: { label?: string; id: string };
  next?: { label?: string; id: string };
} => {
  const currentId = useLocation().pathname;
  const { navigation } = useCurrentNavigation();

  let prev: { label?: string; id: string } | undefined;
  let next: { label?: string; id: string } | undefined;

  let foundCurrent = false;

  traverseNavigation(navigation, (item) => {
    const itemId =
      item.type === "doc"
        ? joinUrl(item.path || "")
        : item.type === "category" && item.path
          ? joinUrl(item.path)
          : undefined;

    if (!itemId) return undefined;

    if (foundCurrent) {
      next = item.label ? { id: itemId, label: item.label } : { id: itemId };
      return true;
    }

    if (currentId === itemId) {
      foundCurrent = true;
    } else {
      prev = item.label ? { id: itemId, label: item.label } : { id: itemId };
    }
    return undefined;
  });

  return {
    ...(next && { next }),
    ...(prev && { prev }),
  };
};

export const navigationListItem = cva(
  "relative flex items-center gap-2 px-(--padding-nav-item) my-0.5 py-1.5 rounded-lg hover:bg-accent tabular-nums",
  {
    defaultVariants: {
      isActive: false,
    },
    variants: {
      isActive: {
        false: "text-foreground/80",
        true: "bg-accent font-medium",
      },
      isMuted: {
        false: "",
        true: "text-foreground/30",
      },
      isPending: {
        false: "",
        true: "bg-accent animate-pulse",
      },
    },
  },
);
