"use client";

import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from "@/components/ui/command";
import { siteConfig } from "@/config/site";
import { cn } from "@/lib/utils";
import { type DialogProps } from "@radix-ui/react-dialog";
import { LaptopIcon, MoonIcon, SunIcon } from "@radix-ui/react-icons";
import * as VisuallyHidden from "@radix-ui/react-visually-hidden";
import { useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import * as React from "react";

import { DialogDescription, DialogTitle } from "./ui/dialog";

export function CommandMenu({ ...props }: DialogProps) {
  const router = useRouter();
  const [open, setOpen] = React.useState(false);
  const { setTheme } = useTheme();

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || e.key === "/") {
        if (
          (e.target instanceof HTMLElement && e.target.isContentEditable) ||
          e.target instanceof HTMLInputElement ||
          e.target instanceof HTMLTextAreaElement ||
          e.target instanceof HTMLSelectElement
        ) {
          return;
        }

        e.preventDefault();
        setOpen((open) => !open);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const runCommand = React.useCallback((command: () => unknown) => {
    setOpen(false);
    command();
  }, []);

  return (
    <div className="w-full">
      <Button
        className={cn(
          "h-8 w-full justify-between gap-2 rounded-lg bg-muted/50 text-base font-normal text-muted-foreground"
        )}
        onClick={() => setOpen(true)}
        variant="outline"
        {...props}
      >
        <span className="inline-flex">Type a command or search...</span>
        <kbd className="pointer-events-none flex h-5 select-none items-center gap-1 rounded border bg-muted p-2 font-mono text-[10px] font-medium">
          <span className="text-xs">⌘</span>
          <span>K</span>
        </kbd>
      </Button>
      <CommandDialog onOpenChange={setOpen} open={open}>
        <VisuallyHidden.Root>
          <DialogTitle />
        </VisuallyHidden.Root>
        <DialogDescription />
        <CommandInput placeholder="Type a command or search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          {Object.entries(siteConfig.links).map(([heading, links]) => (
            <CommandGroup
              heading={heading.charAt(0).toUpperCase() + heading.slice(1)}
              key={heading}
            >
              {links.map((navItem) => (
                <CommandItem
                  className="flex gap-2"
                  key={navItem.href}
                  onClick={() => {
                    runCommand(() => router.push(navItem.href as string));
                  }}
                  onSelect={() => {
                    runCommand(() => router.push(navItem.href as string));
                  }}
                  value={navItem.title}
                >
                  <navItem.Icon className="h-5 w-5" />
                  <span>{navItem.title}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          ))}
          <CommandSeparator />
          <CommandGroup heading="Theme">
            {["light", "dark", "system"].map((theme) => (
              <CommandItem
                className="flex gap-2"
                key={theme}
                onSelect={() => runCommand(() => setTheme(theme))}
              >
                {theme === "light" && <SunIcon className="h-5 w-5" />}
                {theme === "dark" && <MoonIcon className="h-5 w-5" />}
                {theme === "system" && <LaptopIcon className="h-5 w-5" />}
                <span>{theme.charAt(0).toUpperCase() + theme.slice(1)}</span>
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </CommandDialog>
    </div>
  );
}
