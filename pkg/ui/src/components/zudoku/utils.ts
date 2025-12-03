import type { ClassValue } from "clsx";
import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function joinUrl(...parts: (string | undefined)[]) {
  return (
    parts
      .filter((part): part is string => Boolean(part))
      .join("/")
      .replace(/\/+/g, "/")
      .replace(/\/$/, "") || "/"
  );
}

export function shouldShowItem(item: unknown, condition?: unknown): boolean {
  if (!item) return false;
  if (condition === undefined) return true;
  if (typeof condition === "boolean") return condition;
  if (typeof condition === "function") return condition(item);
  return true;
}

export function normalizeUrl(url: string): string {
  if (!url) return "/";
  if (url.startsWith("http://") || url.startsWith("https://")) return url;
  return url.startsWith("/") ? url : `/${url}`;
}
