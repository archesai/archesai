import type { LucideIcon } from "lucide-react";
import {
  AlertTriangleIcon,
  InfoIcon,
  LightbulbIcon,
  ShieldAlertIcon,
} from "lucide-react";
import type { JSX, ReactNode } from "react";

import { cn } from "#lib/utils";

const stylesMap = {
  caution: {
    bg: "bg-yellow-100/60 dark:bg-yellow-400/10",
    border: "border-yellow-400 dark:border-yellow-400/25",
    Icon: AlertTriangleIcon as LucideIcon,
    iconColor: "text-yellow-500 dark:text-yellow-300",
    textColor: "text-yellow-700 dark:text-yellow-200",
    titleColor: "text-yellow-600 dark:text-yellow-300",
  },
  danger: {
    bg: "bg-rose-50 dark:bg-rose-950/40",
    border: "border-rose-400 dark:border-rose-800",
    Icon: ShieldAlertIcon as LucideIcon,
    iconColor: "text-rose-400 dark:text-rose-300",
    textColor: "text-rose-700 dark:text-rose-100",
    titleColor: "text-rose-800 dark:text-rose-300",
  },
  info: {
    bg: "bg-blue-50 dark:bg-blue-950/40",
    border: "border-blue-400 dark:border-blue-900/60",
    Icon: InfoIcon as LucideIcon,
    iconColor: "text-blue-400 dark:text-blue-200",
    textColor: "text-blue-600 dark:text-blue-100",
    titleColor: "text-blue-700 dark:text-blue-200",
  },
  note: {
    bg: "bg-gray-100 dark:bg-zinc-800/50",
    border: "border-gray-300 dark:border-zinc-800",
    Icon: InfoIcon as LucideIcon,
    iconColor: "text-gray-600 dark:text-zinc-300",
    textColor: "text-gray-600 dark:text-zinc-300",
    titleColor: "text-gray-600 dark:text-zinc-300",
  },
  tip: {
    bg: "bg-green-200/25 dark:bg-green-950/70",
    border: "border-green-500 dark:border-green-800",
    Icon: LightbulbIcon as LucideIcon,
    iconColor: "text-green-600 dark:text-green-200",
    textColor: "text-green-600 dark:text-green-50",
    titleColor: "text-green-700 dark:text-green-200",
  },
} as const;

interface CalloutProps {
  children: ReactNode;
  className?: string;
  icon?: boolean;
  title?: string;
  type: keyof typeof stylesMap;
}

export const Callout = ({
  children,
  className,
  icon = true,
  title,
  type,
}: CalloutProps): JSX.Element => {
  const { bg, border, Icon, iconColor, textColor, titleColor } =
    stylesMap[type];

  return (
    <div
      className={cn(
        "not-prose my-2 rounded-md border p-4 text-md",
        icon &&
          "grid grid-cols-[min-content_1fr] grid-rows-[fit-content_1fr] items-baseline gap-x-4 gap-y-2",
        !icon && title && "flex flex-col gap-2",
        "[&_a]:underline [&_a]:decoration-current [&_a]:decoration-from-font [&_a]:underline-offset-4 hover:[&_a]:decoration-1",
        "[&_.code-block-wrapper]:border",
        "[&_ol]:list-decimal [&_ul>li]:my-1 [&_ul>li]:ps-1 [&_ul]:list-disc [&_ul]:ps-4",
        icon && title && "items-center",
        border,
        bg,
        className,
      )}
    >
      {icon && (
        <Icon
          aria-hidden="true"
          className={cn(!title ? "translate-y-1" : "align-middle", iconColor)}
          size={20}
        />
      )}
      {title && <h3 className={cn("font-medium", titleColor)}>{title}</h3>}
      <div
        className={cn(
          icon && "col-start-2",
          !title && icon && "row-start-1",
          textColor,
          "overflow-x-auto",
        )}
      >
        {children}
      </div>
    </div>
  );
};
