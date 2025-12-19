import { useTheme } from "next-themes";
import type { JSX } from "react";

import { MoonIcon, SunIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";

export function ThemeToggle(): JSX.Element {
  const { setTheme, theme } = useTheme();

  return (
    <Button
      className="group/toggle extend-touch-target"
      onClick={() => {
        setTheme(theme === "dark" ? "light" : "dark");
      }}
      size="sm"
      title="Toggle theme"
      variant="ghost"
    >
      <SunIcon className="size-4 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <MoonIcon className="absolute size-4 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
