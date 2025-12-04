import { cn } from "../utils";

export const ColorMap = {
  blue: "bg-sky-400 dark:bg-sky-800",
  gray: "bg-gray-400 dark:bg-gray-600",
  green: "bg-green-400 dark:bg-green-800",
  indigo: "bg-indigo-400 dark:bg-indigo-600",
  outline: "border border-border rounded-md text-foreground",
  purple: "bg-purple-400 dark:bg-purple-600",
  red: "bg-red-400 dark:bg-red-800",
  yellow: "bg-yellow-400 dark:bg-yellow-800",
};

export const ColorMapInvert = {
  blue: "text-sky-400 dark:text-sky-600",
  gray: "text-gray-400 dark:text-gray-600",
  green: "text-green-500 dark:text-green-600",
  indigo: "text-indigo-400 dark:text-indigo-600",
  outline: "",
  purple: "text-purple-400 dark:text-purple-600",
  red: "text-red-400 dark:text-red-600",
  yellow: "text-yellow-400 dark:text-yellow-600",
};

export const NavigationBadge = ({
  color,
  label,
  className,
  invert,
}: {
  color: keyof typeof ColorMap;
  label: string;
  className?: string;
  invert?: boolean;
}) => {
  return (
    <span
      className={cn(
        "flex h-full items-center rounded-sm text-center font-bold text-[0.65rem] text-background uppercase leading-5 transition-opacity duration-200 dark:text-zinc-50",
        color === "outline" ? "px-3" : "mt-0.5 px-1",
        invert ? ColorMapInvert[color] : ColorMap[color],
        className,
      )}
    >
      {label}
    </span>
  );
};
