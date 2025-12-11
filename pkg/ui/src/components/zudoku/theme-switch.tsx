import { MoonIcon, SunIcon } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "#components/shadcn/button";
import { ClientOnly } from "./client-only";
import { cn } from "./utils";

export const ThemeSwitch = () => {
  const { resolvedTheme, setTheme } = useTheme();

  return (
    <ClientOnly
      fallback={
        <Button
          size="icon"
          variant="ghost"
        />
      }
    >
      <button
        aria-label={
          resolvedTheme === "dark"
            ? "Switch to light mode"
            : "Switch to dark mode"
        }
        className="group flex gap-0.5 rounded-full border p-0.5"
        onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}
        type="button"
      >
        <div
          className={cn(
            "rounded-full border border-transparent p-0.5 transition-all [&>svg>circle]:fill-transparent [&>svg>circle]:transition-colors [&>svg>circle]:duration-500 [&>svg>path]:transition-transform [&>svg>path]:duration-200",
            resolvedTheme === "light" && "border-border bg-muted",
            resolvedTheme === "dark" &&
              "group-hover:rotate-[15deg] group-hover:[&>svg>path]:-translate-x-[1px] group-hover:[&>svg>path]:-translate-y-[1px] group-hover:[&>svg>path]:scale-110",
          )}
        >
          <SunIcon size={16} />
        </div>
        <div
          className={cn(
            "transform-gpu rounded-full border border-transparent p-0.5 transition-transform duration-500",
            resolvedTheme === "dark" &&
              "border-border bg-muted [&>svg>path]:fill-white [&>svg>path]:stroke-transparent",
            resolvedTheme === "light" &&
              "group-hover:rotate-[-10deg] [&>svg>path]:stroke-currentColor",
          )}
        >
          <MoonIcon size={16} />
        </div>
      </button>
    </ClientOnly>
  );
};
